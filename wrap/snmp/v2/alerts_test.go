package v2_test

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	of "github.com/cisco-cx/of/pkg/v2"
	of_snmp "github.com/cisco-cx/of/pkg/v2/snmp"
	logger "github.com/cisco-cx/of/wrap/logrus/v2"
	snmp "github.com/cisco-cx/of/wrap/snmp/v2"
	uuid "github.com/cisco-cx/of/wrap/uuid/v2"
	yaml "github.com/cisco-cx/of/wrap/yaml/v2"
)

// Enforce AlertGenerator Interface
func TestAlertsInterface(t *testing.T) {
	var _ of_snmp.AlertGenerator = &snmp.Alerter{}
}

// Test Alerts firing.
func TestAlertFire(t *testing.T) {

	ag := newAlerter(t)

	// All possible alerts for given configs and trapVars.
	alerts, err := ag.Alert([]string{"epc"})
	require.NoError(t, err)

	expectedAlert := []of.Alert{
		of.Alert{
			Labels: map[string]string{
				"alert_name":        "starCard",
				"alert_severity":    "error",
				"source_address":    "192.168.1.28",
				"source_hostname":   "localhost",
				"star_slot_num":     "14",
				"subsystem":         "epc",
				"vendor":            "cisco",
				"alert_fingerprint": "1d8540881d5a50ae",
				"event_id":          "9dcc77fc-dda5-4edf-a683-64f2589036d6",
			},
			Annotations: map[string]string{
				"alert_name":                "starCard",
				"alert_severity":            "error",
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

	require.Equal(t, expectedAlert, alerts)

}

// Test Alerts clearing.
func TestAlertClear(t *testing.T) {

	ag := newAlerter(t)

	// All possible alerts for given configs and trapVars.
	alerts, err := ag.Alert([]string{"nso"})
	require.NoError(t, err)

	expectedAlert := of.Alert{
		Labels: map[string]string{
			"alert_severity":    "error",
			"alertname":         "nsoPackageLoadFailure",
			"source_address":    "nso1.example.org",
			"source_hostname":   "nso1.example.org",
			"subsystem":         "nso",
			"alert_fingerprint": "17da3a8a1001fb2d",
			"vendor":            "cisco",
			"event_id":          "9dcc77fc-dda5-4edf-a683-64f2589036d6",
		},
		Annotations: map[string]string{
			"alert_severity":            "error",
			"event_type":                "clear",
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
		EndsAt:       time.Now().Format(of.AMTimeFormat),
		GeneratorURL: "http://www.oid-info.com/get/1.3.6.1.4.1.8164.2.44",
	}

	// EndsAt is time.Now, so individually matching other components.
	require.Len(t, alerts, 1)
	require.Equal(t, expectedAlert.Labels, alerts[0].Labels)
	require.Equal(t, expectedAlert.Annotations, alerts[0].Annotations)
	require.Equal(t, expectedAlert.GeneratorURL, alerts[0].GeneratorURL)
	require.NotEmpty(t, alerts[0].EndsAt)

}

// Preparing Alert generator.
func newAlerter(t *testing.T) *snmp.Alerter {
	// Prepare snmp.V2Config
	r := strings.NewReader(YamlContent)
	cfg := yaml.Configs{}
	err := cfg.Decode(r)
	require.NoError(t, err)
	configs := of_snmp.V2Config(cfg)

	//Preparing logger
	l := logger.New()

	mr := mibRegistry(t)

	return &snmp.Alerter{
		Log:      l,
		Configs:  &configs,
		Receipts: TrapReceipts(),
		Value:    snmp.NewValue(trapVars(), mr),
		Mr:       mr,
		U:        &uuid.FixedUUID{},
	}
}
