package v2_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"

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
	require.Len(as.t, *alerts, 4)

	nonEmpty := 0
	for idx, alert := range *alerts {
		if alert.EndsAt.IsZero() == false {
			nonEmpty += 1
			(*alerts)[idx].EndsAt = time.Time{}
		}
	}

	require.Equal(as.t, nonEmpty, 3)
	expectedAlerts := []of.Alert{
		of.Alert{
			Labels: map[string]string{
				"alert_name":        "starCard",
				"alert_severity":    "error",
				"source_address":    "192.168.1.28",
				"source_hostname":   "localhost",
				"star_slot_num":     "14",
				"subsystem":         "epc",
				"vendor":            "cisco",
				"alert_fingerprint": "5dd1df6eff3119f4",
				"alert_oid":         ".1.3.6.1.4.1.8164.1.2.1.1.1",
			},
			Annotations: map[string]string{
				"event_id":                  "9dcc77fc-dda5-4edf-a683-64f2589036d6",
				"event_oid":                 ".1.3.6.1.4.1.8164.1.2.1.1.1",
				"event_type":                "error",
				"source_address":            "192.168.1.28",
				"source_hostname":           "localhost",
				"event_description":         "",
				"event_filebeat_timestamp":  "2019-04-26T03:46:57.941Z",
				"event_name":                "unknown",
				"event_oid_string":          "",
				"event_rawtext":             "SNMPTRAP timestamp=[2019-04-26T03:46:57Z] hostname=[localhost] address=[UDP/IPv6: [::1]:48381] pdu_security=[TRAP2, SNMP v3, user user-sha-aes128, context ] vars[.1.3.6.1.2.1.1.3.0 = Timeticks: (123) 0:00:01.23\t.1.3.6.1.6.3.1.1.4.1.0 = OID: .1.3.6.1.6.3.1.1.5.1\t.1.3.6.1.4.1.65000.1.1.1.1.1 = STRING: \"foo\"\t.1.3.6.1.4.1.65000.1.1.1.1.1 = STRING: \"bar\"]",
				"event_snmptrapd_timestamp": "2019-04-26T03:46:57Z",
				"event_vars_json":           "[{\"oid\":\".1.3.6.1.6.1.1.1.4.1\",\"oid_string\":\"1.3.6.1.6.1.1.1.4.1\",\"oid_uri\":\"http://www.oid-info.com/get/1.3.6.1.6.1.1.1.4.1\",\"type\":\"\",\"value\":\".1.3.6.1.4.1.8164.1.2.1.1.1\"},{\"oid\":\".1.3.6.1.4.1.8164.1.2.1.1.1\",\"oid_string\":\"1.3.6.1.4.1.8164.1.2.1.1.1\",\"oid_uri\":\"http://www.oid-info.com/get/1.3.6.1.4.1.8164.1.2.1.1.1\",\"type\":\"\",\"value\":\"14\"},{\"oid\":\".1.3.6.1.4.1.24961.2.103.1.1.5.1.2\",\"oid_string\":\"1.3.6.1.4.1.24961.2.103.1.1.5.1.2\",\"oid_uri\":\"http://www.oid-info.com/get/1.3.6.1.4.1.24961.2.103.1.1.5.1.2\",\"type\":\"\",\"value\":\"package-load-failure\"},{\"oid\":\".1.3.6.1.2.1.1.3.0\",\"oid_string\":\"1.3.6.1.2.1.1.3.0\",\"oid_uri\":\"http://www.oid-info.com/get/1.3.6.1.2.1.1.3.0\",\"type\":\"Timeticks\",\"value\":\"(123) 0:00:01.23\"},{\"oid\":\".1.3.6.1.6.3.1.1.4.1\",\"oid_string\":\"1.3.6.1.6.3.1.1.4.1\",\"oid_uri\":\"http://www.oid-info.com/get/1.3.6.1.6.3.1.1.4.1\",\"type\":\"OID\",\"value\":\".1.3.6.1.4.1.8164.2.13\"},{\"oid\":\".1.3.6.1.6.3.1.1.4.1.0\",\"oid_string\":\"1.3.6.1.6.3.1.1.4.1.0\",\"oid_uri\":\"http://www.oid-info.com/get/1.3.6.1.6.3.1.1.4.1.0\",\"type\":\"OID\",\"value\":\".1.3.6.1.4.1.8164.2.44\"},{\"oid\":\".1.3.6.1.4.1.8164.2.44\",\"oid_string\":\"1.3.6.1.4.1.8164.2.44\",\"oid_uri\":\"http://www.oid-info.com/get/1.3.6.1.4.1.8164.2.44\",\"type\":\"STRING\",\"value\":\"foo\"},{\"oid\":\".1.3.6.1.6.3.1.1.4.1.1\",\"oid_string\":\"1.3.6.1.6.3.1.1.4.1.1\",\"oid_uri\":\"http://www.oid-info.com/get/1.3.6.1.6.3.1.1.4.1.1\",\"type\":\"OID\",\"value\":\".1.3.6.1.4.1.8164.2.45\"},{\"oid\":\".1.3.6.1.4.1.8164.2.45\",\"oid_string\":\"1.3.6.1.4.1.8164.2.45\",\"oid_uri\":\"http://www.oid-info.com/get/1.3.6.1.4.1.8164.2.45\",\"type\":\"OID\",\"value\":\".1.3.6.1.4.1.65000.1.1.1.1.1\"},{\"oid\":\".1.3.6.1.4.1.65000.1.1.1.1.1\",\"oid_string\":\"1.3.6.1.4.1.65000.1.1.1.1.1\",\"oid_uri\":\"http://www.oid-info.com/get/1.3.6.1.4.1.65000.1.1.1.1.1\",\"type\":\"STRING\",\"value\":\"bar\"}]",
			},
			GeneratorURL: "http://www.oid-info.com/get/1.3.6.1.4.1.8164.2.44",
		},
	}

	var err error
	expectedAlerts[0].StartsAt, err = time.Parse(time.RFC3339, "2019-04-26T03:46:57Z")
	require.NoError(as.t, err)

	expectedClearAlertTemplate := of.Alert{
		Labels: map[string]string{
			"alert_severity":    "error",
			"alertname":         "nsoPackageLoadFailure",
			"source_address":    "nso1.example.org",
			"source_hostname":   "nso1.example.org",
			"subsystem":         "nso",
			"alert_fingerprint": "ec92aefbceeb3cd4",
			"vendor":            "cisco",
		},
		Annotations: map[string]string{
			"event_type":                "clear",
			"event_id":                  "9dcc77fc-dda5-4edf-a683-64f2589036d6",
			"source_address":            "nso1.example.org",
			"source_hostname":           "nso1.example.org",
			"event_description":         "",
			"event_filebeat_timestamp":  "2019-04-26T03:46:57.941Z",
			"event_name":                "unknown",
			"event_oid_string":          "",
			"event_rawtext":             "SNMPTRAP timestamp=[2019-04-26T03:46:57Z] hostname=[localhost] address=[UDP/IPv6: [::1]:48381] pdu_security=[TRAP2, SNMP v3, user user-sha-aes128, context ] vars[.1.3.6.1.2.1.1.3.0 = Timeticks: (123) 0:00:01.23\t.1.3.6.1.6.3.1.1.4.1.0 = OID: .1.3.6.1.6.3.1.1.5.1\t.1.3.6.1.4.1.65000.1.1.1.1.1 = STRING: \"foo\"\t.1.3.6.1.4.1.65000.1.1.1.1.1 = STRING: \"bar\"]",
			"event_snmptrapd_timestamp": "2019-04-26T03:46:57Z",
			"event_vars_json":           "[{\"oid\":\".1.3.6.1.6.1.1.1.4.1\",\"oid_string\":\"1.3.6.1.6.1.1.1.4.1\",\"oid_uri\":\"http://www.oid-info.com/get/1.3.6.1.6.1.1.1.4.1\",\"type\":\"\",\"value\":\".1.3.6.1.4.1.8164.1.2.1.1.1\"},{\"oid\":\".1.3.6.1.4.1.8164.1.2.1.1.1\",\"oid_string\":\"1.3.6.1.4.1.8164.1.2.1.1.1\",\"oid_uri\":\"http://www.oid-info.com/get/1.3.6.1.4.1.8164.1.2.1.1.1\",\"type\":\"\",\"value\":\"14\"},{\"oid\":\".1.3.6.1.4.1.24961.2.103.1.1.5.1.2\",\"oid_string\":\"1.3.6.1.4.1.24961.2.103.1.1.5.1.2\",\"oid_uri\":\"http://www.oid-info.com/get/1.3.6.1.4.1.24961.2.103.1.1.5.1.2\",\"type\":\"\",\"value\":\"package-load-failure\"},{\"oid\":\".1.3.6.1.2.1.1.3.0\",\"oid_string\":\"1.3.6.1.2.1.1.3.0\",\"oid_uri\":\"http://www.oid-info.com/get/1.3.6.1.2.1.1.3.0\",\"type\":\"Timeticks\",\"value\":\"(123) 0:00:01.23\"},{\"oid\":\".1.3.6.1.6.3.1.1.4.1\",\"oid_string\":\"1.3.6.1.6.3.1.1.4.1\",\"oid_uri\":\"http://www.oid-info.com/get/1.3.6.1.6.3.1.1.4.1\",\"type\":\"OID\",\"value\":\".1.3.6.1.4.1.8164.2.13\"},{\"oid\":\".1.3.6.1.6.3.1.1.4.1.0\",\"oid_string\":\"1.3.6.1.6.3.1.1.4.1.0\",\"oid_uri\":\"http://www.oid-info.com/get/1.3.6.1.6.3.1.1.4.1.0\",\"type\":\"OID\",\"value\":\".1.3.6.1.4.1.8164.2.44\"},{\"oid\":\".1.3.6.1.4.1.8164.2.44\",\"oid_string\":\"1.3.6.1.4.1.8164.2.44\",\"oid_uri\":\"http://www.oid-info.com/get/1.3.6.1.4.1.8164.2.44\",\"type\":\"STRING\",\"value\":\"foo\"},{\"oid\":\".1.3.6.1.6.3.1.1.4.1.1\",\"oid_string\":\"1.3.6.1.6.3.1.1.4.1.1\",\"oid_uri\":\"http://www.oid-info.com/get/1.3.6.1.6.3.1.1.4.1.1\",\"type\":\"OID\",\"value\":\".1.3.6.1.4.1.8164.2.45\"},{\"oid\":\".1.3.6.1.4.1.8164.2.45\",\"oid_string\":\"1.3.6.1.4.1.8164.2.45\",\"oid_uri\":\"http://www.oid-info.com/get/1.3.6.1.4.1.8164.2.45\",\"type\":\"OID\",\"value\":\".1.3.6.1.4.1.65000.1.1.1.1.1\"},{\"oid\":\".1.3.6.1.4.1.65000.1.1.1.1.1\",\"oid_string\":\"1.3.6.1.4.1.65000.1.1.1.1.1\",\"oid_uri\":\"http://www.oid-info.com/get/1.3.6.1.4.1.65000.1.1.1.1.1\",\"type\":\"STRING\",\"value\":\"bar\"}]",
			"event_oid":                 ".1.3.6.1.4.1.8164.2.44",
		},
		EndsAt:       time.Time{},
		GeneratorURL: "http://www.oid-info.com/get/1.3.6.1.4.1.8164.2.44",
	}

	for _, val := range []string{
		".1.3.6.1.4.1.24961.2.103.2.0.3",
		".1.3.6.1.4.1.24961.2.103.2.0.4",
		".1.3.6.1.4.1.24961.2.103.2.0.5",
	} {
		expectedClearAlertTemplate.Labels["alert_oid"] = val
		expectedClearAlertTemplate.StartsAt, err = time.Parse(time.RFC3339, "2019-04-26T03:46:57Z")
		require.NoError(as.t, err)
		expectedAlerts = append(expectedAlerts, expectedClearAlertTemplate)
	}

	require.ElementsMatch(as.t, expectedAlerts, *alerts)
	return nil
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

	// Send SNMP traps to server.
	dataBytes, err := json.Marshal(TrapEvents())
	require.NoError(t, err)

	data := bytes.NewBuffer(dataBytes)
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
	r := strings.NewReader(YamlContent)
	configs := yaml.Configs{}
	err := configs.Decode(r)
	require.NoError(t, err)

	v2Config := of_snmp.V2Config(configs)

	// Prepare MIBS registry
	mr := mibRegistry(t)

	// Prepare lookup.
	lookup := snmp.Lookup{Configs: v2Config, MR: mr, Log: l}

	err = lookup.Build()
	require.NoError(t, err)

	u := uuid.FixedUUID{}

	cntr, cntrVec := snmp.InitCounters(namespace, l)

	ag := snmp.Alerter{
		Log:     l,
		Configs: &v2Config,
		MR:      mr,
		U:       &u,
		Cntr:    cntr,
		CntrVec: cntrVec,
	}

	// INIT SNMP service.
	s := &snmp.Service{
		Writer:  herodot.New(l),
		Log:     l,
		MR:      mr,
		Configs: &v2Config,
		U:       &u,
		As:      &testAlertService{t: t},
		Lookup:  &lookup,
		Alerter: &ag,
		Cntr:    cntr,
		CntrVec: cntrVec,
	}
	return s
}
