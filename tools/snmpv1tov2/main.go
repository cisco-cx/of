package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
	of_snmpv1 "github.com/cisco-cx/of/pkg/v1/snmp"
	of_snmpv2 "github.com/cisco-cx/of/pkg/v2/snmp"
)

func main() {
	if len(os.Args) != 3 {
		usage()
	}

	inputFile := os.Args[1]
	outputFile := os.Args[2]

	// Open input file
	r, err := os.Open(inputFile)
	if err != nil {
		fmt.Printf("Failed to open input file, %s", err.Error())
		return
	}
	defer r.Close()

	// Parse old config.
	cfg := of_snmpv1.AlertsConfig{}
	err = yaml.NewDecoder(r).Decode(&cfg)
	if err != nil {
		fmt.Printf("Failed to decode input file, %s", err.Error())
		return
	}

	// Prepare new config.
	V2Config := of_snmpv2.V2Config{}
	newCfg := of_snmpv2.Config{}
	newCfg.Defaults.SourceType = of_snmpv2.HostType
	newCfg.Alerts = make([]of_snmpv2.Alert, 0)
	for name, alert := range cfg.Alerts {
		newAlert := newAlert(name, alert)
		newCfg.Alerts = append(newCfg.Alerts, newAlert)
	}

	// Write new config
	V2Config[filepath.Base(inputFile)] = newCfg

	w, err := os.Create(outputFile)
	if err != nil {
		fmt.Printf("Failed to open output file, %s", err.Error())
		return
	}
	defer w.Close()
	err = yaml.NewEncoder(w).Encode(&V2Config)
	if err != nil {
		fmt.Printf("Failed to encode output file, %s", err.Error())
		return
	}
}

// Prepare new alert.
func newAlert(alertName string, alert *of_snmpv1.AlertEntry) of_snmpv2.Alert {
	// Init alert
	newAlert := of_snmpv2.Alert{}
	newAlert.LabelMods = make([]of_snmpv2.Mod, 0)
	newAlert.AnnotationMods = make([]of_snmpv2.Mod, 0)

	// alert name in config.
	newAlert.Name = alertName

	// Mod to set alert name into label
	alertNameMod := of_snmpv2.Mod{
		Type:  of_snmpv2.Set,
		Key:   "alertname",
		Value: alertName,
	}
	newAlert.LabelMods = append(newAlert.LabelMods, alertNameMod)

	// Mod to set alert severity into label
	alertSeverityMod := of_snmpv2.Mod{
		Type:  of_snmpv2.Set,
		Key:   "alert_severity",
		Value: alert.AlertSeverity,
	}
	newAlert.LabelMods = append(newAlert.LabelMods, alertSeverityMod)

	// Get values that should match to fire/clear an alert, and its corresponding event name.
	firingValues, firingEventMap := valuesAndEventMap(alert.ErrorEvents)
	clearingValues, clearingEventMap := valuesAndEventMap(alert.ClearEvents)

	clearEvents, _ := json.Marshal(alert.ClearEvents)

	// Firing Selects
	if len(firingValues) != 0 {
		firingSelect := of_snmpv2.Select{
			Type:   of_snmpv2.Equals,
			Oid:    of_snmpv2.SNMPTrapOID,
			As:     of_snmpv2.Value,
			Values: firingValues,
			AnnotationMods: []of_snmpv2.Mod{
				of_snmpv2.Mod{
					Type:  of_snmpv2.Copy,
					Oid:   of_snmpv2.SNMPTrapOID,
					As:    of_snmpv2.Value,
					ToKey: "event_name",
					Map:   firingEventMap,
				},
				of_snmpv2.Mod{
					Type:  of_snmpv2.Set,
					Key:   "compatible_clear_events",
					Value: string(clearEvents),
				},
			},
		}

		newAlert.Firing = map[string][]of_snmpv2.Select{
			"select": []of_snmpv2.Select{firingSelect},
		}
	}

	// Clearing Selects
	if len(clearingValues) != 0 {
		clearingSelect := of_snmpv2.Select{
			Type:   of_snmpv2.Equals,
			Oid:    of_snmpv2.SNMPTrapOID,
			As:     of_snmpv2.Value,
			Values: clearingValues,
			AnnotationMods: []of_snmpv2.Mod{
				of_snmpv2.Mod{
					Type:  of_snmpv2.Copy,
					Oid:   of_snmpv2.SNMPTrapOID,
					As:    of_snmpv2.Value,
					ToKey: "event_name",
					Map:   clearingEventMap,
				},
			},
		}
		newAlert.Clearing = map[string][]of_snmpv2.Select{
			"select": []of_snmpv2.Select{clearingSelect},
		}
	}

	return newAlert
}

// Prepare values that should match to fire/clear an alert, and its corresponding event name.
func valuesAndEventMap(events map[string]of_snmpv1.Event) ([]string, map[string]string) {
	values := make([]string, 0)
	eventMap := make(map[string]string)
	for valueOID, eventName := range events {
		values = append(values, valueOID)
		eventMap[valueOID] = eventName.EventName
	}

	return values, eventMap
}

func usage() {
	fmt.Printf("%s <input_file <output_file>\n", filepath.Base(os.Args[0]))
	os.Exit(1)
}
