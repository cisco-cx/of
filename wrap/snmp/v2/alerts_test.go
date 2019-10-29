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
	yaml "github.com/cisco-cx/of/wrap/yaml/v2"
)

// Enforce AlertGenerator Interface
func TestAlertsInterface(t *testing.T) {
	var _ of_snmp.AlertGenerator = &snmp.Alerter{}
}

// Test Alerts firing.
func TestAlertFire(t *testing.T) {

	// Prepare snmp.V2Config
	r := strings.NewReader(YamlContent)
	cfg := yaml.Configs{}
	err := cfg.Decode(r)
	require.NoError(t, err)
	configs := of_snmp.V2Config(cfg)

	//Preparing logger
	l := logger.New()

	// Preparing Alert generator.
	ag := snmp.Alerter{
		Log:     l,
		Configs: &configs,
		Source:  TrapSource(),
		Value:   snmp.NewValue(trapVars(), mibRegistry(t)),
	}

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
				"alert_fingerprint": "39a7842eabe0437a",
			},
			Annotations: map[string]string{
				"alert_name":      "starCard",
				"alert_severity":  "error",
				"event_oid":       ".1.3.6.1.4.1.8164.1.2.1.1.1",
				"event_type":      "error",
				"source_address":  "192.168.1.28",
				"source_hostname": "localhost",
			},
			GeneratorURL: "http://www.oid-info.com/get/1.3.6.1.4.1.8164.2.13",
		},
	}

	require.Equal(t, expectedAlert, alerts)

}

// Test Alerts clearing.
func TestAlertClear(t *testing.T) {

	// Prepare snmp.V2Config
	r := strings.NewReader(YamlContent)
	cfg := yaml.Configs{}
	err := cfg.Decode(r)
	require.NoError(t, err)
	configs := of_snmp.V2Config(cfg)

	//Preparing logger
	l := logger.New()

	// Preparing Alert generator.
	ag := snmp.Alerter{
		Log:     l,
		Configs: &configs,
		Source:  TrapSource(),
		Value:   snmp.NewValue(trapVars(), mibRegistry(t)),
	}

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
			"alert_fingerprint": "362a7c9e679338f1",
			"vendor":            "cisco",
		},
		Annotations: map[string]string{
			"alert_severity":  "error",
			"event_type":      "clear",
			"source_address":  "nso1.example.org",
			"source_hostname": "nso1.example.org",
		},
		EndsAt:       time.Now().Format(of.AMTimeFormat),
		GeneratorURL: "http://www.oid-info.com/get/1.3.6.1.4.1.8164.2.13",
	}

	// EndsAt is time.Now, so individually matching other components.
	require.Len(t, alerts, 1)
	require.Equal(t, expectedAlert.Labels, alerts[0].Labels)
	require.Equal(t, expectedAlert.Annotations, alerts[0].Annotations)
	require.Equal(t, expectedAlert.GeneratorURL, alerts[0].GeneratorURL)
	require.NotEmpty(t, alerts[0].EndsAt)

}
