package v2

import (
	"encoding/json"

	of "github.com/cisco-cx/of/pkg/v2"
	of_snmp "github.com/cisco-cx/of/pkg/v2/snmp"
	am "github.com/cisco-cx/of/wrap/alertmanager/v2"
	concatenator "github.com/cisco-cx/of/wrap/concatenator/v2"
	herodot "github.com/cisco-cx/of/wrap/herodot/v2"
	logger "github.com/cisco-cx/of/wrap/logrus/v2"
	mib_registry "github.com/cisco-cx/of/wrap/mib/v2"
	mibs "github.com/cisco-cx/of/wrap/mib/v2"
	prometheus "github.com/cisco-cx/of/wrap/prometheus/client_golang/v2"
	uuid "github.com/cisco-cx/of/wrap/uuid/v2"
	yaml "github.com/cisco-cx/of/wrap/yaml/v2"
)

const (
	// Counters names.
	clearingEventCount    = "clearing_alert_count"
	unknownClusterIPCount = "unknown_cluster_ip_count"
	eventsReceivedCount   = "events_received_count"
	eventsProcessedCount  = "events_processed_count"

	//CounterVec names.
	alertsGeneratedCount    = "alerts_generated_count"
	alertsNotGeneratedCount = "alerts_not_generated_count"
	unknownAlertsCount      = "unknown_alerts_count"
	HandlerRestarted        = "handler_restarted"
	alertsGenerationFailed  = "alerts_generation_failed_count"
)

type Service struct {
	Writer     of.Writer
	Log        *logger.Logger
	Configs    *of_snmp.V2Config
	MR         *mibs.MIBRegistry
	U          of.UUIDGen
	Lookup     of_snmp.Lookup
	As         of.Notifier
	Cntr       map[string]*prometheus.Counter
	CntrVec    map[string]*prometheus.CounterVec
	SNMPConfig *of.SNMPConfig
}

func NewService(l *logger.Logger, cfg *of.SNMPConfig, cntr map[string]*prometheus.Counter, cntrVec map[string]*prometheus.CounterVec) (*Service, error) {

	// Concatenate configs files.
	c := concatenator.Files{
		Path: cfg.ConfigDir,
		Ext:  "yaml",
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
		DryRun:    cfg.DryRun,
	}

	// Prepare lookup.
	lookup := Lookup{Configs: v2Config, MR: mr, Log: l}
	err = lookup.Build()
	if err != nil {
		l.WithError(err).Errorf("Failed to build lookup.")
		return nil, err
	}

	u := uuid.UUID{}

	// INIT SNMP service.
	s := &Service{
		Writer:     herodot.New(l),
		Log:        l,
		MR:         mr,
		Configs:    &v2Config,
		U:          &u,
		As:         &as,
		Lookup:     &lookup,
		Cntr:       cntr,
		CntrVec:    cntrVec,
		SNMPConfig: cfg,
	}
	return s, nil
}

// Search lookup to find configs that match Trap Vars values.
func (s Service) lookupConfigs(events []*of.PostableEvent) [][]string {
	configs := make([][]string, len(events))
	for idx, event := range events {
		s.Cntr[eventsReceivedCount].Incr()
		s.Log.Tracef("Event[%d] %+v", idx, event)
		cfgs, err := s.Lookup.Find(&event.Document.Receipts.Snmptrapd.Vars)
		if err != nil {
			s.Log.WithError(err).Errorf("Lookup failed.")
			continue
		}
		configs[idx] = cfgs
		if len(cfgs) == 0 {
			valueLookup := NewValue(&event.Document.Receipts.Snmptrapd.Vars, s.MR)
			trapV, _ := valueLookup.Value(of_snmp.SNMPTrapOID)
			trapVStrShort, _ := valueLookup.ValueStrShort(of_snmp.SNMPTrapOID)
			s.Log.WithFields(map[string]interface{}{
				"vars":             event.Document.Receipts.Snmptrapd.Vars,
				"source":           event.Document.Receipts.Snmptrapd.Source,
				"event":            event,
				"SNMPTrapOIDValue": trapV,
				"SNMPTrapOIDName":  trapVStrShort,
			}).Debugf("No match found at lookup")
		}
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

	alerter := Alerter{
		Log:            s.Log,
		Configs:        s.Configs,
		MR:             s.MR,
		U:              s.U,
		Cntr:           s.Cntr,
		CntrVec:        s.CntrVec,
		LogUnknown:     s.SNMPConfig.LogUnknown,
		ForwardUnknown: s.SNMPConfig.ForwardUnknown,
	}

	var alerts []of.Alert
	for index, event := range events {
		snmptrapd := event.Document.Receipts.Snmptrapd
		alerter.Receipts = &event.Document.Receipts
		alerter.Value = NewValue(&snmptrapd.Vars, s.MR)

		trapV, _ := alerter.Value.Value(of_snmp.SNMPTrapOID)
		s.Log.WithFields(map[string]interface{}{
			"index":            index,
			"timestamp":        snmptrapd.Timestamp,
			"source":           snmptrapd.Source,
			"SNMPTrapOIDValue": trapV,
		}).Infof("Processing event")
		if len(configs[index]) != 0 {
			alerts = append(alerts, alerter.Alert(configs[index])...)
		} else {
			alerts = append(alerts, alerter.Unknown("lookup")...)
		}

		s.Cntr[eventsProcessedCount].Incr()
	}
	var fireAlerts []of.Alert
	var clearAlerts []of.Alert
	for _, alert := range alerts {
		if alert.EndsAt.IsZero() {
			clearAlerts = append(clearAlerts, alert)
		} else {
			fireAlerts = append(fireAlerts, alert)
		}
	}
	s.Log.Infof("Generated %d alerts firing : %d, clearing : %d", len(alerts), len(fireAlerts), len(clearAlerts))
	err := s.As.Notify(&fireAlerts)
	if err != nil {
		s.Log.WithError(err).Errorf("Failed to publish firing alert(s) for received event")
		s.Writer.WriteCode(w, r, 503, nil)
		return
	}
	err = s.As.Notify(&clearAlerts)
	if err != nil {
		s.Log.WithError(err).Errorf("Failed to publish clearing alert(s) for received event")
		s.Writer.WriteCode(w, r, 503, nil)
		return
	}
	s.Writer.WriteCode(w, r, 200, nil)
}

// Create counters..
func InitCounters(namespace string, log *logger.Logger) (map[string]*prometheus.Counter, map[string]*prometheus.CounterVec) {
	if namespace == "" {
		log.Fatalf("Counters namespace cannot be empty.")
	}
	cntr := map[string]*prometheus.Counter{
		clearingEventCount: &prometheus.Counter{Namespace: namespace, Name: clearingEventCount,
			Help: "Number of unique clear events generated."},
		unknownClusterIPCount: &prometheus.Counter{Namespace: namespace, Name: unknownClusterIPCount,
			Help: "Number of times we got events for a cluster, where the IP is not in the cluster list."},
		eventsReceivedCount: &prometheus.Counter{Namespace: namespace, Name: eventsReceivedCount,
			Help: "Number of SNMP trap events receieved."},
		eventsProcessedCount: &prometheus.Counter{Namespace: namespace, Name: eventsProcessedCount,
			Help: "Number of SNMP trap events processed by the handler."},
	}

	// Init counters
	for name, c := range cntr {
		err := c.Create()
		if err != nil {
			log.WithError(err).Fatalf("Failed to init counter, %s", name)
		}
	}

	// Represents the details needed to init a counter vector.
	type vectorInfo struct {
		vector *prometheus.CounterVec
		labels []string
	}

	// Available vectors
	vis := []vectorInfo{
		vectorInfo{
			vector: &prometheus.CounterVec{
				Namespace: namespace,
				Name:      alertsGeneratedCount,
				Help:      "Number of times we generated an alert object for sending to AlertManager.",
			},
			labels: []string{"alertType", "alert_oid"},
		},

		vectorInfo{
			vector: &prometheus.CounterVec{
				Namespace: namespace,
				Name:      alertsNotGeneratedCount,
				Help:      "Number of times alert were not generated for configs that matched our lookup.",
			},
			labels: []string{"level", "alert_oid"},
		},
		vectorInfo{
			vector: &prometheus.CounterVec{
				Namespace: namespace,
				Name:      unknownAlertsCount,
				Help:      "Number of times we encountered unknown SNMP traps.",
			},
			labels: []string{"level", "alert_oid"},
		},
		vectorInfo{
			vector: &prometheus.CounterVec{
				Namespace: namespace,
				Name:      HandlerRestarted,
				Help:      "Number of times handler was restarted.",
			},
			labels: []string{"op_type"},
		},
		vectorInfo{
			vector: &prometheus.CounterVec{
				Namespace: namespace,
				Name:      alertsGenerationFailed,
				Help:      "Number of times alert generation failed.",
			},
			labels: []string{"alertType", "alert_oid"},
		},
	}

	cntrVec := make(map[string]*prometheus.CounterVec)
	for _, vi := range vis {
		err := vi.vector.Create(vi.labels)
		if err != nil {
			log.WithError(err).Fatalf("Failed to init counterVec, %s", vi.vector.Name)
		}
		cntrVec[vi.vector.Name] = vi.vector
	}
	return cntr, cntrVec
}
