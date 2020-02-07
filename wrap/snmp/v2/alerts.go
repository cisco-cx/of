package v2

import (
	"encoding/json"
	"strings"
	"time"

	prommodel "github.com/prometheus/common/model"
	of "github.com/cisco-cx/of/pkg/v2"
	of_snmp "github.com/cisco-cx/of/pkg/v2/snmp"
	logger "github.com/cisco-cx/of/wrap/logrus/v2"
	mibs "github.com/cisco-cx/of/wrap/mib/v2"
	prometheus "github.com/cisco-cx/of/wrap/prometheus/client_golang/v2"
)

// Implements of_snmp.AlertGenerator
type Alerter struct {
	Log            *logger.Logger
	Configs        *of_snmp.V2Config
	Receipts       *of.Receipts
	Value          *Value
	MR             *mibs.MIBRegistry
	U              of.UUIDGen
	Cntr           map[string]*prometheus.Counter
	CntrVec        map[string]*prometheus.CounterVec
	LogUnknown     bool
	ForwardUnknown bool
}

// Iterate through configs in configNames and generate all possible Alerts.
func (a *Alerter) Alert(cfgNames []string) []of.Alert {

	if a.Cntr == nil || len(a.Cntr) == 0 {
		a.Log.Panicf("Counters have not been initiated.")
	}

	// Fixed annotionations for this set of Trap vars.
	fixedAnnotations := a.fixedAnnotations()
	var allAlerts = make([]of.Alert, 0)
	for _, cfgName := range cfgNames {

		var cfg of_snmp.Config
		var ok bool

		// Check if config exists
		if cfg, ok = (*a.Configs)[cfgName]; ok == false {
			a.Log.WithError(of.ErrConfigNotFound).Errorf("")
			continue
		}

		trapV, err := a.Value.Value(of_snmp.SNMPTrapOID)
		if err != nil {
			a.Log.WithError(err).Errorf("Failed to get SNMPTrapOID's value.")
		}
		trapVStrShort, err := a.Value.ValueStrShort(of_snmp.SNMPTrapOID)
		if err != nil {
			a.Log.WithError(err).Errorf("Failed to get SNMPTrapOID value's short name.")
		}

		// Identify device.
		a.Log.WithFields(map[string]interface{}{
			"vars":        a.Receipts.Snmptrapd.Vars,
			"PduSecurity": a.Receipts.Snmptrapd.PduSecurity,
			"config":      cfgName,
		}).Debugf("Trying to identify device.")

		if a.deviceIdentified(cfg.Defaults.DeviceIdentifiers) == false {
			a.Log.WithFields(map[string]interface{}{
				"vars":        a.Receipts.Snmptrapd.Vars,
				"PduSecurity": a.Receipts.Snmptrapd.PduSecurity,
				"config":      cfgName,
			}).Debugf("Config not applicable for device.")
			continue
		}

		// To check if any alert was generated for `cfgName`.
		var alertMatchedConfig bool = false
		// Check through of_snmp.Config.Alerts to find a match.
		for aNum, alertCfg := range cfg.Alerts {

			// To check if any alert was generated for `alertCfg`.
			var alertMatchedAlertCfg bool = false

			// Check if alert is enabled.
			a.Log.Debugf("Checking Alert no. %d (%v) in config: %s", aNum, alertCfg.Name, cfgName)
			if a.enabled(cfg.Defaults.Enabled, alertCfg.Enabled) == false {
				// Printing alert no., since alert name can be nil.
				a.Log.Debugf("Alert no. %d (%v), not enabled in config: %s", aNum, alertCfg.Name, cfgName)
				continue
			}

			// Check if trap Vars have any alert matching firing conditions.
			fAlert, err := a.matchAlerts(cfg, alertCfg, of_snmp.Firing, fixedAnnotations)
			if err == nil {
				alertMatchedAlertCfg = true
				alertMatchedConfig = true

				fAlert.Annotations[string(of_snmp.EventTypeText)] = string(of_snmp.Firing)
				// Setting `alert_oid` as the value of of_snmp.SNMPTrapOID
				fAlert.Labels["alert_oid"] = fAlert.Annotations["event_oid"]
				a.CntrVec[alertsGeneratedCount].Incr(map[string]string{
					"alertType": "firing",
					"alert_oid": fAlert.Labels["alert_oid"],
				})
				a.StartsAt(&fAlert)
				a.EndsAt(cfg.Defaults.EndsAt, alertCfg.EndsAt, &fAlert)

				// Finger print the alert.
				fingerprint := a.fingerprint(fAlert)
				fAlert.Labels[of_snmp.FingerprintText] = fingerprint

				allAlerts = append(allAlerts, fAlert)
				a.Log.WithFields(map[string]interface{}{
					"alertType":   "firing",
					"labels":      fAlert.Labels,
					"annotations": fAlert.Annotations,
					"startsAt":    fAlert.StartsAt,
					"endsAt":      fAlert.EndsAt,
					"vars":        a.Receipts.Snmptrapd.Vars,
					"source":      a.Receipts.Snmptrapd.Source,
					"config":      cfgName,
				}).Debugf("Generated alerts")
				alertJson, _ := json.Marshal(fAlert)
				a.Log.Tracef("alert_json : %+v", string(alertJson))
			}

			// Check if trap Vars have any alerts matching clearing conditions.
			cAlert, err := a.matchAlerts(cfg, alertCfg, of_snmp.Clearing, fixedAnnotations)
			if err == nil {
				alertMatchedAlertCfg = true
				alertMatchedConfig = true

				// Add end time to clearing alerts.
				cAlert.Annotations[string(of_snmp.EventTypeText)] = string(of_snmp.Clearing)
				cAlert.EndsAt = time.Now().UTC()
				a.StartsAt(&cAlert)
				a.EndsAt(cfg.Defaults.EndsAt, alertCfg.EndsAt, &cAlert)

				a.Cntr[clearingEventCount].Incr()
				// For `selects` under firing.
				for _, s := range alertCfg.Firing["select"] {
					// Add each OID under values as `alert_oid`
					for _, v := range s.Values {
						// Setting `alert_oid` to clear for all known firing values.
						cAlert.Labels["alert_oid"] = v
						a.CntrVec[alertsGeneratedCount].Incr(map[string]string{
							"alertType": "clearing",
							"alert_oid": v,
						})

						// Finger print the alert.
						fingerprint := a.fingerprint(cAlert)
						cAlert.Labels[of_snmp.FingerprintText] = fingerprint

						allAlerts = append(allAlerts, cAlert)
						a.Log.WithFields(map[string]interface{}{
							"alertType":   "clearing",
							"labels":      cAlert.Labels,
							"annotations": cAlert.Annotations,
							"startsAt":    cAlert.StartsAt,
							"endsAt":      cAlert.EndsAt,
							"vars":        a.Receipts.Snmptrapd.Vars,
							"source":      a.Receipts.Snmptrapd.Source,
							"config":      cfgName,
						}).Debugf("Generated alerts")
						alertJson, _ := json.Marshal(cAlert)
						a.Log.Tracef("alert_json : %+v", string(alertJson))
					}
				}
			}

			// Alert not generated for `alertCfg`
			if alertMatchedAlertCfg == false {
				a.Log.WithFields(map[string]interface{}{
					"alertCfg":         alertCfg,
					"alertIndex":       aNum,
					"vars":             a.Receipts.Snmptrapd.Vars,
					"source":           a.Receipts.Snmptrapd.Source,
					"config":           cfgName,
					"SNMPTrapOIDValue": trapV,
					"SNMPTrapOIDName":  trapVStrShort,
				}).Debugf("No match found for alertCfg.")

				a.CntrVec[alertsNotGeneratedCount].Incr(map[string]string{
					"level":     "alert",
					"alert_oid": trapV,
				})
			}
		}

		// Alert not generated for `cfgName`
		if alertMatchedConfig == false {
			a.Log.WithFields(map[string]interface{}{
				"config":           cfgName,
				"vars":             a.Receipts.Snmptrapd.Vars,
				"source":           a.Receipts.Snmptrapd.Source,
				"SNMPTrapOIDValue": trapV,
				"SNMPTrapOIDName":  trapVStrShort,
			}).Debugf("No match found for config.")
			a.CntrVec[alertsNotGeneratedCount].Incr(map[string]string{
				"level":     "config",
				"alert_oid": trapV,
			})
			allAlerts = append(allAlerts, a.Unknown("config")...)
		}
	}
	return allAlerts
}

func (a *Alerter) StartsAt(alert *of.Alert) {
	t, err := time.Parse(time.RFC3339, a.Receipts.Snmptrapd.Timestamp)
	if err != nil {
		a.Log.WithError(err).Errorf("Failed to parse time from %s", a.Receipts.Snmptrapd.Timestamp)
		return
	}

	alert.StartsAt = t
}

func (a *Alerter) EndsAt(defEndsAt int, alertEndsAt int, alert *of.Alert) {
	if defEndsAt != 0 {
		alert.EndsAt = time.Now().UTC().Add(time.Duration(defEndsAt) * time.Minute)
	}

	if alertEndsAt != 0 {
		alert.EndsAt = time.Now().UTC().Add(time.Duration(alertEndsAt) * time.Minute)
	}
}

// Create alert for unknown SNMP trap.
func (a *Alerter) Unknown(level string) []of.Alert {

	trapV, err := a.Value.Value(of_snmp.SNMPTrapOID)
	if err != nil {
		a.Log.WithError(err).Errorf("Failed to get SNMPTrapOID's value.")
	}
	trapVStrShort, err := a.Value.ValueStrShort(of_snmp.SNMPTrapOID)
	if err != nil {
		a.Log.WithError(err).Errorf("Failed to get SNMPTrapOID value's short name.")
	}

	a.CntrVec[unknownAlertsCount].Incr(map[string]string{
		"level":     level,
		"alert_oid": trapV,
	})

	info := a.Log.WithFields(map[string]interface{}{
		"level":            level,
		"vars":             a.Receipts.Snmptrapd.Vars,
		"source":           a.Receipts.Snmptrapd.Source,
		"SNMPTrapOIDValue": trapV,
		"SNMPTrapOIDName":  trapVStrShort,
	})

	if a.LogUnknown == true {
		info.Infof("Unknown alert.")
	} else {
		info.Debugf("Unknown alert.")
	}

	if a.ForwardUnknown == false {
		return []of.Alert{}
	}

	var alert = of.Alert{}

	// Init new Alert
	alert.Labels = make(map[string]string)
	alert.Annotations = make(map[string]string)

	// Fixed annotionations for this set of Trap vars.
	fixedAnnotations := a.fixedAnnotations()
	// Apply fixed annotionations.
	for k, v := range fixedAnnotations {
		alert.Annotations[k] = v
	}

	alert.Labels["alertname"] = "unknownSnmpTrap"
	alert.Labels["alert_oid"] = alert.Annotations["event_oid"]
	alert.Labels["source_address"] = a.Receipts.Snmptrapd.Source.Address
	alert.Labels["source_hostname"] = a.Receipts.Snmptrapd.Source.Hostname

	alert.Labels[of_snmp.FingerprintText] = a.fingerprint(alert)

	return []of.Alert{alert}
}

// match trapVars with alerts in config.
func (a *Alerter) matchAlerts(cfg of_snmp.Config, alertCfg of_snmp.Alert, alertType of_snmp.EventType, fixedAnnotations map[string]string) (of.Alert, error) {

	var alert = of.Alert{}

	var selects []of_snmp.Select
	switch alertType {
	case of_snmp.Firing:
		selects = alertCfg.Firing["select"]
	case of_snmp.Clearing:
		selects = alertCfg.Clearing["select"]
	default:
		return alert, of.ErrUnknownEventType
	}

	// Check if trap Vars have any alerts matching select conditions.
	matched, err := a.selected(selects)
	if err != nil {
		a.Log.WithError(err).Errorf("Error while trying to match alert.")
		return alert, err
	}

	if matched == false {
		a.Log.WithError(of.ErrNoMatch).WithField("alertType", alertType).Debugf("")
		return alert, of.ErrNoMatch
	}

	a.Log.WithField("alertType", alertType).Debugf("Alert matched.")

	// Prepare alert since the selects have matched.

	// Init new Alert
	alert.Labels = make(map[string]string)
	alert.Annotations = make(map[string]string)

	// Apply fixed annotionations.
	for k, v := range fixedAnnotations {
		alert.Annotations[k] = v
	}

	// Preparing base alert.
	err = a.prepareBaseAlert(&alert, &cfg)
	if err != nil {
		a.Log.WithError(err).Errorf("Error while preparing base alert.")
		return alert, err
	}

	// Apply alert specific labels.
	alert.Annotations["event_id"] = a.U.UUID()
	err = a.applyMod(&alert.Labels, alertCfg.LabelMods)
	if err != nil {
		a.Log.WithError(err).Errorf("Error while applying alert mods to labels.")
		return alert, err
	}

	// Apply alert specific annotations.
	err = a.applyMod(&alert.Annotations, alertCfg.AnnotationMods)
	if err != nil {
		a.Log.WithError(err).Errorf("Error while applying alert mods to annotations.")
		return alert, err
	}

	// Generator URL prefix.
	alert.GeneratorURL = string(a.generatorURLPrefix(cfg.Defaults.GeneratorUrlPrefix, alertCfg.GeneratorUrlPrefix))
	SNMPTrapOIDValue, err := a.Value.Value(of_snmp.SNMPTrapOID)
	if err == nil {
		alert.GeneratorURL += strings.TrimPrefix(SNMPTrapOIDValue, ".")
	}

	// Apply select specfic changes.
	for _, sel := range selects {

		// Apply select specific annotations.
		err = a.applyMod(&alert.Annotations, sel.AnnotationMods)
		if err != nil {
			a.Log.WithError(err).Errorf("Error while applying alert mods to annotations.")
			return alert, err
		}

	}

	return alert, nil
}

// Check if give selects are applicable.
func (a *Alerter) selected(selects []of_snmp.Select) (bool, error) {
	// If none selects are mentioned, then match should fail.
	if len(selects) == 0 {
		return false, nil
	}

	allSelectsMatched := true
	for _, sel := range selects {
		// Find value based on the of_snmp.As type.
		resolvedValue, err := a.Value.ValueAs(sel.Oid, sel.As)
		if err != nil {
			a.Log.WithError(err).Errorf("Failed to resolve %s as %s.", sel.Oid, sel.As)
			return false, err
		}

		// Check if value is present in of_snmp.Config.Alerts[x].Select[y].Values
		valueFound := false
		for _, v := range sel.Values {
			if v == resolvedValue {
				valueFound = true
				break
			}
		}

		// If value is not found in of_snmp.Config.Alerts[x].Select[y].Values, stop checking for other selects.
		if valueFound == false {
			allSelectsMatched = false
			break
		}

	}

	return allSelectsMatched, nil
}

// Prepares the base alert based on keys under of_snmp.Config.Defaults
func (a *Alerter) prepareBaseAlert(alert *of.Alert, cfg *of_snmp.Config) error {

	// Update source info.
	var found = false
	if cfg.Defaults.SourceType == of_snmp.ClusterType {
		// Iterate through clusters to check if IP matches with available IP.
		for clusterName, cluster := range cfg.Defaults.Clusters {
			for _, ip := range cluster.SourceAddresses {
				if ip == a.Receipts.Snmptrapd.Source.Address {
					a.Log.Debugf("Found cluster name %s, for source IP : %s", clusterName, ip)
					found = true
					alert.Labels["source_address"] = clusterName
					alert.Labels["source_hostname"] = clusterName
					alert.Annotations["source_address"] = clusterName
					alert.Annotations["source_hostname"] = clusterName
					break
				}
			}
		}

	}

	// If no cluster is found or host type is not cluster.
	if found == false {
		a.Log.Debugf("Setting default source info for IP : %s", a.Receipts.Snmptrapd.Source.Address)
		a.Cntr[unknownClusterIPCount].Incr()
		a.updateSource(alert)
	}

	// Apply default mods to Labels
	err := a.applyMod(&(alert.Labels), cfg.Defaults.LabelMods)
	if err != nil {
		a.Log.WithError(err).Errorf("Failed to apply default mods to labels.")
		return err
	}

	// Apply default mods to Annotations
	err = a.applyMod(&(alert.Annotations), cfg.Defaults.AnnotationMods)
	if err != nil {
		a.Log.WithError(err).Errorf("Failed to apply default mods to annotations.")
		return err
	}

	return nil
}

// Annotations that don't change for a SNMP trap event.
func (a *Alerter) fixedAnnotations() map[string]string {

	enrichedVars := make([]map[string]string, len(a.Receipts.Snmptrapd.Vars))
	var oid, eventStrOid, eventDesc string
	for i, v := range a.Receipts.Snmptrapd.Vars {

		enrichedVar := make(map[string]string)
		enrichedVar["oid"] = v.Oid
		enrichedVar["type"] = v.Type
		enrichedVar["value"] = v.Value

		varOid := v.Oid[1:]
		object := a.MR.MIB(varOid)
		if object != nil {
			enrichedVar["name"] = object.Name
			enrichedVar["description"] = object.Description
			enrichedVar["units"] = object.Units
		}
		enrichedVar["oid_string"] = a.MR.String(varOid)
		enrichedVar["oid_uri"] = "http://www.oid-info.com/get/" + varOid
		enrichedVars[i] = enrichedVar
		if v.Oid == of_snmp.SNMPTrapOID {
			oid = v.Value
			eventOid := oid[1:]
			eventObj := a.MR.MIB(eventOid)
			if eventObj != nil {
				eventDesc = eventObj.Description
				eventStrOid = a.MR.String(eventOid)
			}
		}
	}

	enrichedVarsJson, _ := json.Marshal(enrichedVars)

	fixedAnnotations := map[string]string{
		"event_filebeat_timestamp":  a.Receipts.Filebeat.Timestamp,
		"event_name":                "unknown",
		"event_oid":                 oid,
		"event_oid_string":          eventStrOid,
		"event_rawtext":             a.Receipts.Filebeat.Message,
		"event_snmptrapd_timestamp": a.Receipts.Snmptrapd.Timestamp,
		"event_type":                "unknown",
		"event_vars_json":           string(enrichedVarsJson),
		"event_description":         eventDesc,
	}

	return fixedAnnotations
}

func (a *Alerter) applyMod(mapPtr *map[string]string, mods []of_snmp.Mod) error {

	// Init modifier
	m := Modifier{
		V: a.Value,
	}

	a.Log.WithField("value", a.Value).Tracef("Mod values")

	// Apply mods
	m.Map = mapPtr
	err := m.Apply(mods)
	if err != nil {
		return err
	}
	return nil
}

// Update source info based on of.TrapSource
func (a *Alerter) updateSource(alert *of.Alert) {
	alert.Labels["source_address"] = a.Receipts.Snmptrapd.Source.Address
	alert.Labels["source_hostname"] = a.Receipts.Snmptrapd.Source.Hostname
	alert.Annotations["source_address"] = a.Receipts.Snmptrapd.Source.Address
	alert.Annotations["source_hostname"] = a.Receipts.Snmptrapd.Source.Hostname
}

// Overwrite with alert specific prefix if available.
func (a *Alerter) generatorURLPrefix(defPrefix of_snmp.URLPrefix, alertPrefix of_snmp.URLPrefix) of_snmp.URLPrefix {
	if defPrefix == "" && alertPrefix == "" {
		return ""
	}

	// If alertPrefix is defined, use it.
	if alertPrefix != "" {
		return alertPrefix
	}

	return defPrefix
}

// Identify device based on string in PduSecurity.
// If no identifier are present in the config, consider the config for alerts.
func (a *Alerter) deviceIdentified(identifiers []string) bool {
	if len(identifiers) == 0 {
		return true
	}

	for _, identifier := range identifiers {
		if strings.Contains(a.Receipts.Snmptrapd.PduSecurity, identifier) == true {
			return true
		}
	}
	return false
}

// Decide if alert should be sent or not based the of_snmp.Config.Defaults.Enabled and of_snmp.Config.Alerts[name].Enabled
//
// defaults.enabled 	alerts[n].enabled 	State
// Undefined 			Undefined 			Enabled
// Undefined 			false 				Disabled
// Undefined 			true 				Enabled
// false 				Undefined 			Disabled
// false 				false 				Disabled
// false 				true 				Disabled
// true 				Undefined 			Enabled
// true 				false 				Disabled
// true 				true 				Enabled
func (a *Alerter) enabled(defEnabled *bool, alertEnabled *bool) bool {
	// If both are undefined.
	if defEnabled == nil && alertEnabled == nil {
		return true
	}

	// If default is undefined, but alertEnabled is defined.
	if defEnabled == nil {
		return *alertEnabled
	}

	// If default is false.
	if *defEnabled == false {
		return false
	}

	// If alertEnabled is undefined.
	if alertEnabled == nil {
		return *defEnabled
	}

	// If defEnabled is true.
	return *alertEnabled
}

// Fingerprint the alert.
func (a *Alerter) fingerprint(al of.Alert) string {
	labels := make(prommodel.LabelSet)
	for k, v := range al.Labels {
		labels[prommodel.LabelName(k)] = prommodel.LabelValue(v)
	}
	return labels.Fingerprint().String()
}
