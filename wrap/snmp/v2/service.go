package v2

import (
	"encoding/json"
	"fmt"

	of "github.com/cisco-cx/of/pkg/v2"
	of_snmp "github.com/cisco-cx/of/pkg/v2/snmp"
	am "github.com/cisco-cx/of/wrap/alertmanager/v2"
	concatenator "github.com/cisco-cx/of/wrap/concatenator/v2"
	herodot "github.com/cisco-cx/of/wrap/herodot/v2"
	logger "github.com/cisco-cx/of/wrap/logrus/v2"
	mib_registry "github.com/cisco-cx/of/wrap/mib/v2"
	mibs "github.com/cisco-cx/of/wrap/mib/v2"
	uuid "github.com/cisco-cx/of/wrap/uuid/v2"
	yaml "github.com/cisco-cx/of/wrap/yaml/v2"
)

type Service struct {
	Writer  of.Writer
	Log     *logger.Logger
	Configs *of_snmp.V2Config
	MR      *mibs.MIBRegistry
	U       of.UUIDGen
	Lookup  of_snmp.Lookup
	As      of.Notifier
	CN      string
}

func NewService(l *logger.Logger, cfg *of.SNMPConfig) (*Service, error) {

	// Concatenate configs files.
	c := concatenator.Files{
		Path: cfg.ConfigDir,
	}

	r, err := c.Concat()
	if err != nil {
		l.WithError(err).Errorf("Failed to concat config files in %s.", cfg.ConfigDir)
		return nil, err
	}

	// Decode configs files.
	configs := yaml.Configs{}
	err = configs.Decode(r)
	if err != nil {
		l.WithError(err).Errorf("Failed to decode config files in %s.", cfg.ConfigDir)
		return nil, err
	}

	v2Config := of_snmp.V2Config(configs)

	// Prepare MIBS registry
	mr := mib_registry.New()

	readerMIB := &mib_registry.MIBHandler{
		MapMIB: make(map[string]of.MIB),
	}

	if cfg.CacheFile != "none" {
		err = readerMIB.LoadCacheFromFile(cfg.CacheFile)
		if err != nil {
			l.WithError(err).Fatalf("Failed to load MIBs from cache.")
		}
	} else {
		if cfg.SNMPMibsDir == "" {
			l.Fatalf("Failed to load MIBs, no cache path or SNMP MIBs Dir.")
		}
		err = readerMIB.LoadJSONFromDir(cfg.SNMPMibsDir)
		if err != nil {
			l.WithError(err).Fatalf("Failed to load MIBs from MIBS dir.")
		}
	}

	err = mr.Load(readerMIB.MapMIB)
	if err != nil {
		l.WithError(err).Fatalf("Failed to load MIBs.")
	}

	// Setup alert service.
	as := am.AlertService{
		Version:   cfg.Version,
		AmURL:     cfg.AMAddress,
		Throttle:  cfg.Throttle,
		PostTime:  cfg.PostTime,
		SleepTime: cfg.SleepTime,
		SendTime:  cfg.SendTime,
		Log:       l,
	}

	// Prepare lookup.
	lookup := Lookup{Configs: v2Config, MR: mr}
	err = lookup.Build()
	if err != nil {
		l.WithError(err).Errorf("Failed to build lookup.")
		return nil, err
	}

	// INIT SNMP service.
	s := &Service{
		Writer:  herodot.New(l),
		Log:     l,
		MR:      mr,
		Configs: &v2Config,
		U:       &uuid.UUID{},
		As:      &as,
		Lookup:  &lookup,
		CN:      cfg.Application,
	}
	return s, nil
}

// Search lookup to find configs that match Trap Vars values.
func (s Service) lookupConfigs(events []*of.PostableEvent) [][]string {
	configs := make([][]string, len(events))
	for idx, event := range events {
		cfgs, err := s.Lookup.Find(&event.Document.Receipts.Snmptrapd.Vars)
		if err != nil {
			s.Log.WithError(err).Errorf("Lookup failed.")
			continue
		}
		configs[idx] = cfgs
	}
	return configs
}

// HTTP handler func to recieve SNMP traps.
func (s Service) AlertHandler(w of.ResponseWriter, r of.Request) {

	// Decode http.Body into of.PostableEvent
	var events []*of.PostableEvent
	if err := json.NewDecoder(r.Body).Decode(&events); err != nil {
		s.Log.WithError(err).Errorf("Failed to decode events.")
		s.Writer.WriteError(w, r, herodot.ErrBadRequest.WithError(err.Error()))
	}
	s.Log.Infof("Received %d events.", len(events))

	configs := s.lookupConfigs(events)

	ag := Alerter{
		Log:     s.Log,
		Configs: s.Configs,
		MR:      s.MR,
		U:       s.U,
		CN:      s.CN,
	}

	for index, event := range events {
		snmptrapd := event.Document.Receipts.Snmptrapd
		s.Log.WithFields(map[string]interface{}{"index": index, "timestamp": snmptrapd.Timestamp, "hostname": snmptrapd.Source.Hostname}).Debugf("Processing event")
		if len(configs[index]) == 0 {
			s.Log.Debugf("No config found for index, %d", index)
			continue
		}
		//ag.ReValue
		ag.Receipts = &event.Document.Receipts
		ag.Value = NewValue(&snmptrapd.Vars, s.MR)
		alerts, err := ag.Alert(configs[index])
		if err != nil {
			continue
		}
		s.Log.Infof("Generated %d alerts for event[%s]", len(alerts), fmt.Sprintf("%d", index))
		err = s.As.Notify(&alerts)
		if err != nil {
			s.Log.WithError(err).Errorf("Failed to publish alert(s) for received event")
		}
	}
	s.Writer.WriteCode(w, r, 200, nil)
}
