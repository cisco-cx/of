package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v2"
	of_snmpv2 "github.com/cisco-cx/of/pkg/v2/snmp"
)

func main() {
	if len(os.Args) != 2 {
		usage()
	}

	inputFile := os.Args[1]

	// Open input file
	r, err := os.Open(inputFile)
	if err != nil {
		log.Fatalln("Failed to open input file, ", err.Error())
		return
	}
	defer r.Close()

	// Parse config.
	cfg := of_snmpv2.V2Config{}
	err = yaml.NewDecoder(r).Decode(&cfg)
	if err != nil {
		log.Fatalln("Failed to decode input file, ", err.Error())
		return
	}

	// Generate a list of maps to allow render CSV easier
	var rows []map[string]string = make([]map[string]string, 0)
	for _, cfg := range cfg {
		for _, alert := range cfg.Alerts {
			allMods := append(alert.LabelMods, cfg.Defaults.LabelMods...)
			allMods = append(allMods, cfg.Defaults.AnnotationMods...)
			allAlertNames := searchDefaultAlertNames(&allMods)
			if allAlertNames == nil {
				log.Fatalln("Couldn't find the alert name", inputFile)
			}
			for _, alertName := range *allAlertNames {
				row, err := getRowFromMods(alertName, &allMods)
				if err != nil {
					log.Fatalf("Error filling the row %s, %v\n", row, err)
				}
				rows = append(rows, row)
			}
		}
	}
	renderCSV(&rows)
}

// renderCSV renders the CSV output in the stdout
func renderCSV(rows *[]map[string]string) {
	severityToOSS := map[string]string{
		"critical":      "Critical",
		"error":         "Critical",
		"major":         "Major",
		"minor":         "Minor",
		"notice":        "Minor",
		"informational": "Minor",
		"warning":       "Warning",
		"info":          "Info",
	}
	severityToOSSClasification := map[string]string{
		"critical":      "Outage",
		"error":         "Outage",
		"major":         "Deterioration",
		"minor":         "Deterioration",
		"notice":        "Notification",
		"informational": "Notification",
		"warning":       "Notification",
		"info":          "Normal",
	}
	severityToOSSServiceAffected := map[string]string{
		"critical":      "Yes",
		"error":         "Yes",
		"major":         "Yes",
		"minor":         "No",
		"notice":        "No",
		"informational": "No",
		"warning":       "No",
		"info":          "No",
	}
	severityToOSSIncident := map[string]string{
		"critical":      "Yes",
		"error":         "Yes",
		"major":         "Yes",
		"minor":         "Yes",
		"notice":        "Yes",
		"informational": "No",
		"warning":       "Yes",
		"info":          "No",
	}

	header := []string{
		"Alarm Code",
		"Alarm Name",
		"Labels",
		"EMS",
		"Classification",
		"Service Affected",
		"Category",
		"MO",
		"Default Severity",
		"Alarm Type",
		"Alarming Delay (minutes)",
		"Estimated time to Close(Minutes)",
		"Incident Creation",
		"TT Delay (minutes)",
		"WO Delay (minutes)",
		"Clear Name",
		"Clear Trap",
		"Software Release",
		"Description",
		"Impact",
		"Suggestion",
		"Generic Alarm Name",
		"Probable Cause",
	}

	writer := csv.NewWriter(os.Stdout)
	defer writer.Flush()

	if err := writer.Write(header); err != nil {
		log.Fatalln("error writing csv:", err)
	}
	for _, label := range *rows {
		var labelToString []string
		for key, value := range label {
			if key == "alertname" || key == "alert_severity" || key == "alert_description" {
				continue
			}
			labelToString = append(labelToString, fmt.Sprintf("%s=%s", key, value))
		}
		sort.Strings(labelToString)
		row := []string{
			label["alertname"],
			label["alertname"],
			strings.Join(labelToString, " "),
			"OF",
			severityToOSSClasification[label["alert_severity"]],
			severityToOSSServiceAffected[label["alert_severity"]],
			"-",
			"-",
			severityToOSS[label["alert_severity"]],
			"-",
			"-",
			"-",
			severityToOSSIncident[label["alert_severity"]],
			"-",
			"-",
			"-",
			"Yes",
			"-",
			label["alert_description"],
			"-",
			"-",
			"-",
			"-",
		}
		if err := writer.Write(row); err != nil {
			log.Fatalln("error writing csv:", err)
		}
	}
}

// searchDefaultAlertNames searches for the alertname, first it tries looking in the
// alert, if it is not defined then it gets all the alertnames defined in the default
// values
func searchDefaultAlertNames(mods *[]of_snmpv2.Mod) *[]string {
	var alertNames []string
	for _, mod := range *mods {
		// alertname defined in the alert as a set mod
		if mod.Key == "alertname" {
			alertNames = append(alertNames, mod.Value)
			return &alertNames
		}
		// alertname defined in the default values
		if mod.ToKey == "alertname" {
			for alertName, _ := range mod.Map {
				alertNames = append(alertNames, alertName)
			}
			return &alertNames
		}
	}
	return nil
}

// getRowFromMods returns a map[string]string filled with all mods passed as parameter.
// It's used later for rendering a row in the csv output
func getRowFromMods(alertName string, allMods *[]of_snmpv2.Mod) (map[string]string, error) {
	var row map[string]string = map[string]string{
		"alertname": alertName,
	}
	for _, mod := range *allMods {
		if mod.Type == of_snmpv2.Set {
			row[mod.Key] = mod.Value
		} else if mod.Type == of_snmpv2.Copy && len(mod.Map) == 0 {
			row[mod.ToKey] = "*"
		} else if mod.Type == of_snmpv2.Copy && len(mod.Map[alertName]) != 0 {
			row[mod.ToKey] = mod.Map[alertName]
		} else {
			return nil, errors.New(fmt.Sprintf("Mod %v unsupported", mod))
		}
	}
	return row, nil
}

func usage() {
	fmt.Printf("%s <input_file.yaml>\n", filepath.Base(os.Args[0]))
	os.Exit(1)
}
