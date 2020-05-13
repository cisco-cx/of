package v2_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	of "github.com/cisco-cx/of/pkg/v2"
	of_snmp "github.com/cisco-cx/of/pkg/v2/snmp"
	herodot "github.com/cisco-cx/of/wrap/herodot/v2"
	http "github.com/cisco-cx/of/wrap/http/v2"
	logger "github.com/cisco-cx/of/wrap/logrus/v2"
	snmp "github.com/cisco-cx/of/wrap/snmp/v2"
	uuid "github.com/cisco-cx/of/wrap/uuid/v2"
	yaml "github.com/cisco-cx/of/wrap/yaml/v2"
)

// Represents of.Notifier interface
type testAlertService struct {
	t *testing.T
}

// Test Notify function.
func (as *testAlertService) Notify(alerts *[]of.Alert) error {
	as.verifyNotify(alerts)
	return nil
}

func (as *testAlertService) verifyNotify(alerts *[]of.Alert) {
	expectedLabels := []map[string]string{
		map[string]string{
			"alert_fingerprint": "26c5384f07068e37",
			"alert_oid":         ".1.3.6.1.4.1.8164.2.150",
			"alert_severity":    "major",
			"alertname":         "starTaskFailure",
			"source_address":    "dead::beef",
			"source_hostname":   "test-device-01",
			"subsystem":         "config1",
			"vendor":            "cisco",
		},
		map[string]string{
			"alert_fingerprint": "60dfb17ae277c3c8",
			"alert_oid":         ".1.3.6.1.4.1.8164.2.150",
			"alert_severity":    "major",
			"alertname":         "starTaskRestart",
			"source_address":    "dead::beef",
			"source_hostname":   "test-device-01",
			"subsystem":         "config1",
			"vendor":            "cisco",
		},
		map[string]string{
			"alert_fingerprint": "26c5384f07068e37",
			"alert_oid":         ".1.3.6.1.4.1.8164.2.150",
			"alert_severity":    "major",
			"alertname":         "starTaskFailure",
			"source_address":    "dead::beef",
			"source_hostname":   "test-device-01",
			"subsystem":         "config1",
			"vendor":            "cisco",
		},
		map[string]string{
			"alert_fingerprint": "a8e70c001d09366f",
			"alert_oid":         ".1.3.6.1.4.1.8164.2.151",
			"alert_severity":    "major",
			"alertname":         "starTaskRestart",
			"source_address":    "dead::beef",
			"source_hostname":   "test-device-01",
			"subsystem":         "config1",
			"vendor":            "cisco",
		},
	}

	compareLabels(as.t, alerts, expectedLabels)
}

func TestSNMPService(t *testing.T) {

	// Init SNMP service
	s := initService(t, "testingService")

	// Start server to listen for SNMP traps
	addr := "localhost:24932"

	hc := &of.HTTPConfig{ListenAddress: addr}

	srv := http.NewServer(hc)
	srv.HandleFunc("/", s.AlertHandler)
	err := srv.ListenAndServe()
	require.NoError(t, err)

	docs := []of.Document{}
	err = json.Unmarshal([]byte(StarEvents), &docs)
	require.NoError(t, err)

	events := make([]*of.PostableEvent, len(docs))
	for i, doc := range docs {
		events[i] = &of.PostableEvent{
			Document: doc,
		}
	}

	dataJson, err := json.Marshal(events)
	require.NoError(t, err)

	data := bytes.NewBuffer(dataJson)

	// Send SNMP traps to server.
	c := http.NewClient()
	req, err := http.NewRequest("Post", "http://"+addr, data)
	require.NoError(t, err)
	req.Header.Set("User-Agent", "test")
	req.Header.Set("Content-Type", "application/json")

	res, err := c.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)

	err = srv.Shutdown()
	require.NoError(t, err)
}

// init SNMP service
func initService(t *testing.T, namespace string) *snmp.Service {

	// Logger
	l := logger.New()
	//l.SetLevel("debug")

	// Decode configs files.
	r := strings.NewReader(YamlConfigs)
	configs := yaml.Configs{}
	err := configs.Decode(r)
	require.NoError(t, err)

	v2Config := of_snmp.V2Config(configs)

	// Prepare MIBS registry
	mr := newFakeMibRegistry()

	// Prepare lookup.
	lookup := snmp.Lookup{Configs: v2Config, MR: mr, Log: l}

	err = lookup.Build()
	require.NoError(t, err)

	u := uuid.FixedUUID{}

	cntr, cntrVec := snmp.InitCounters(namespace, l)
	cfg := &of.SNMPConfig{
		LogUnknown:     true,
		ForwardUnknown: true,
	}

	// INIT SNMP service.
	s := &snmp.Service{
		Writer:     herodot.New(l),
		Log:        l,
		MR:         mr,
		Configs:    &v2Config,
		U:          &u,
		As:         &testAlertService{t: t},
		Lookup:     &lookup,
		Cntr:       cntr,
		CntrVec:    cntrVec,
		SNMPConfig: cfg,
	}
	return s
}

func compareLabels(t *testing.T, alerts *[]of.Alert, expectedLabels []map[string]string) {
	alertLabels := make([]map[string]string, len(*alerts))
	for i, a := range *alerts {
		alertLabels[i] = a.Labels
	}
	require.Equal(t, expectedLabels, alertLabels)
}
