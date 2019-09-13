// Copyright © 2019 Cisco Systems, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package v1alpha1

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/mitchellh/mapstructure"
	of "github.com/cisco-cx/of/lib/v1alpha1"
	acigo "github.com/cisco-cx/of/wrap/acigo/v1alpha1"
	alertmanager "github.com/cisco-cx/of/wrap/alertmanager/v1alpha1"
	http "github.com/cisco-cx/of/wrap/http/v1alpha1"
	logger "github.com/cisco-cx/of/wrap/logrus/v1alpha1"
	prometheus "github.com/cisco-cx/of/wrap/prometheus/client_golang/v1alpha1"
	strcase "github.com/cisco-cx/of/wrap/strcase/v1alpha1"
	yaml "github.com/cisco-cx/of/wrap/yaml/v1alpha1"
)

// Alertmanager alert specific constants.
const (
	apicFaultHelpURL        = "https://pubhub.devnetcloud.com/media/apic-mim-ref-411/docs/FAULT-%s.html"
	amAlertFingerprintLabel = "alert_fingerprint"
)

// Counters names.
const (
	amConnectAttemptCount   = "am_connect_attempt_total"
	amConnectErrorCount     = "am_connect_error_count"
	apicConnectAttemptCount = "apic_connect_attempt_total"
	apicConnectErrorCount   = "apic_connect_error_count"
	alertsGeneratedCount    = "alerts_generated_count"
	faultsDroppedCount      = "faults_dropped_count"
	faultsScrapedCount      = "faults_scraped_count"
	faultsMatchedCount      = "faults_matched_count"
	faultsUnmatchedCount    = "faults_unmatched_count"
	notificationCycleCount  = "notification_cycle_count"
)

type Handler struct {
	Config   *of.ACIConfig
	counters map[string]*prometheus.Counter
	server   *http.Server
	Aci      *acigo.ACIService
	Ams      *alertmanager.AlertService
	ac       *yaml.Alerts
	sc       *yaml.Secrets
}

var log = logger.New()

func (h *Handler) Run() {

	h.initCounters()
	srv := of.Server{
		Addr:         h.Config.ListenAddress,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go h.PushAlerts()

	h.server = http.NewServer(srv)

	h.server.HandleFunc("/", func(w of.ResponseWriter, r of.Request) {
		fmt.Fprint(w, h.Config.Version)
	})

	h.server.Handle("/metrics", prometheus.NewHandler())
	err := h.server.ListenAndServe()
	if err != nil {
		log.WithError(err).Fatalf("Failed to listen at %s", h.Config.ListenAddress)
	}
}

func (h *Handler) initCounters() {

	h.counters = map[string]*prometheus.Counter{
		amConnectAttemptCount: &prometheus.Counter{Namespace: h.Config.Application, Name: amConnectAttemptCount,
			Help: "Number of times we tried to connect to AlertManager."},
		amConnectErrorCount: &prometheus.Counter{Namespace: h.Config.Application, Name: amConnectErrorCount,
			Help: "Number of errors encountered while trying to connect to AlertManager."},
		apicConnectAttemptCount: &prometheus.Counter{Namespace: h.Config.Application, Name: apicConnectAttemptCount,
			Help: "Number of times we tried to connect to APIC."},
		apicConnectErrorCount: &prometheus.Counter{Namespace: h.Config.Application, Name: apicConnectErrorCount,
			Help: "Number of errors encountered while trying to connect to APIC."},
		alertsGeneratedCount: &prometheus.Counter{Namespace: h.Config.Application, Name: alertsGeneratedCount,
			Help: "Number of times we generated an alert object for sending to AlertManager."},
		faultsDroppedCount: &prometheus.Counter{Namespace: h.Config.Application, Name: faultsDroppedCount,
			Help: "Number of times we dropped an APIC fault per alerts.yaml."},
		faultsScrapedCount: &prometheus.Counter{Namespace: h.Config.Application, Name: faultsScrapedCount,
			Help: "Number of faults we scraped from APIC."},
		faultsMatchedCount: &prometheus.Counter{Namespace: h.Config.Application, Name: faultsMatchedCount,
			Help: "Number of times we found an alertConfig that mentioned the encountered fault code."},
		faultsUnmatchedCount: &prometheus.Counter{Namespace: h.Config.Application, Name: faultsUnmatchedCount,
			Help: "Number of times we could not find an alertConfig that mentioned the encountered fault code."},
		notificationCycleCount: &prometheus.Counter{Namespace: h.Config.Application, Name: notificationCycleCount,
			Help: "Number of times we tried ran the notification cycle loop."},
	}

	for name, c := range h.counters {
		err := c.Create()
		if err != nil {
			log.WithError(err).Fatalf("Failed to init counter, %s", name)
		}
	}
}

// Pull ACI faults and forward to Alertmanager.
func (h *Handler) PushAlerts() {

	h.ac = &yaml.Alerts{}
	LoadConfig(h.ac, h.Config.AlertsCFGFile)
	h.sc = &yaml.Secrets{}
	LoadConfig(h.sc, h.Config.SecretsCFGFile)

	for {
		h.counters[notificationCycleCount].Incr()
		log.Infof("Running APIC -> AlertManager notification cycle. (cycle-sleep-seconds=%d)\n", h.Config.CycleInterval)

		h.counters[apicConnectAttemptCount].Incr()
		faults, err := h.Aci.Faults()
		if err != nil {
			h.counters[apicConnectErrorCount].Incr()
			log.WithError(err).Errorf("Failed to get faults.")
			return
		}

		alerts, err := h.faultsToAlerts(faults)
		if err != nil {
			log.WithError(err).Errorf("Failed to convert faults to alerts.")
			goto NEXTCYCLE
		}

		// Notify AlertManager. if we have any alerts
		if len(alerts) > 0 {
			h.counters[amConnectAttemptCount].Incr()
			err = h.Ams.Notify(alerts)
			if err != nil {
				h.counters[amConnectErrorCount].Incr()
				log.Errorf("Notification cycle failed. Will retry in %d, %s\n", h.Config.CycleInterval, err.Error())
				goto NEXTCYCLE
			}
		} else {
			log.Errorf("No faults found")
		}

		log.Infof("Notification cycle succeeded. Sleeping for %d seconds.\n", h.Config.CycleInterval)
	NEXTCYCLE:
		time.Sleep(time.Duration(h.Config.CycleInterval) * time.Second)
	}
}

func (h *Handler) faultsToAlerts(faults []of.Map) ([]*of.Alert, error) {
	var alerts []*of.Alert
	for _, mapFault := range faults {

		// Decode fault into struct.
		f := of.ACIFaultRaw{}
		mapstructure.Decode(mapFault, &f)
		fp := acigo.FaultParser{f}
		h.counters[faultsScrapedCount].Incr()

		// If this is in alerts.yaml:dropped_faults, skip it.
		if _, drop := h.ac.DroppedFaults[strings.ToUpper(fp.Fault.Code)]; drop {
			log.Errorf("Dropping fault: %s\n", f)
			h.counters[faultsDroppedCount].Incr()
			continue
		}

		// Create alert boilerplate.
		alert := &of.Alert{
			Annotations: h.Annotations(f),
			Labels:      of.LabelMap{},
		}

		// Get an integer representation of fault severity for numerical comparison.
		s, err := acigo.NewACIFaultSeverityRaw(fp.Fault.Severity)
		if err != nil {
			log.Errorf("Failed to parse severity, %s", err.Error())
			return nil, err
		}
		faultSeverityLevel := s.ID()

		// Parse string date.
		faultCreated, err := fp.Created()
		if err != nil {
			return nil, err
		}
		faultLastTransition, err := fp.LastTransition()
		if err != nil {
			return nil, err
		}

		// Set key values on the alert.
		alert.StartsAt = faultCreated
		alert.GeneratorURL = fmt.Sprintf(apicFaultHelpURL, f.Code)
		alert.Labels["cluster_name"] = of.LabelValue(h.sc.APIC.Cluster.Name)
		sub_id, _ := fp.SubID()
		alert.Labels["sub_id"] = of.LabelValue(sub_id)
		alert.Annotations["source_address"] = h.Config.SourceAddress
		alert.Annotations["source_hostname"] = h.Config.SourceHostname

		// Try to find this fault code in alerts config.
		alertName, newAlertConfig, err := h.GetAlertConfig(f)
		if err == nil && alertName != "" {
			h.counters[faultsMatchedCount].Incr()
			log.Errorf("Found matching fault code in alertsConfig.Alerts.")
			alert.Labels["alertname"] = of.LabelValue(alertName)
			alert.Labels["alert_severity"] = of.LabelValue(newAlertConfig.AlertSeverity)
		} else {
			h.counters[faultsUnmatchedCount].Incr()
			log.Errorf("%s\n", err)
			log.Errorf("Setting default alertname and severity for fault code: %s\n", f.Code)

			// Fall back to the "rule" field in the scraped fault.
			alert.Labels["alertname"] = of.LabelValue(f.Rule)

			// Fall back to the alerts.yaml:defaults.alert_severity.
			alert.Labels["alert_severity"] = of.LabelValue(h.ac.Defaults.AlertSeverity)
		}

		// If fault severity is below fault severity threshhold, consider the fault resolved.
		// That is, use fault["lastTransition"] as "h.EndsAt".
		// Doing this should cause AlertManager to mark any alerts for this fault "resolved" instead of "firing".
		//
		// NOTE: Importantly, this means that other severities than "cleared" may be set status=resolved.

		s, err = acigo.NewACIFaultSeverityRaw(h.ac.APIC.AlertSeverityThreshold)
		if err != nil {
			log.Errorf("Failed to parse severity, %s", err.Error())
			return nil, err
		}

		faultSeverityThresholdLevel := s.ID()
		if faultSeverityLevel < faultSeverityThresholdLevel {
			alert.EndsAt = faultLastTransition
		}

		// Add "alert_fingerprint" label.
		// TODO:
		//h.addFingerprint(alert)

		// Debug sample code.
		b := []byte{}
		b, err = json.Marshal(alert)
		if err != nil {
			return nil, err
		}
		log.Errorf("Alert generated: %s\n", b)

		alerts = append(alerts, alert)
		h.counters[alertsGeneratedCount].Incr()
	}
	return alerts, nil

}

// TODO:
// Finger print alert.
/*
func (h *Handler) addFingerprint(a *of.Alert) {
	a.Labels[amAlertFingerprintLabel] = model.LabelValue(a.Labels.Fingerprint().String())
}
*/

// Wrapper to read a file into an implementation of of.Decoder.
func LoadConfig(cfg of.Decoder, fileName string) {

	f, err := os.Open(fileName)
	if err != nil {
		log.WithError(err).Fatalf("Failed to open file. config-file : %s", fileName)
	}

	err = cfg.Decode(f)
	if err != nil {
		log.WithError(err).Fatalf("Failed to decode alerts config file.")
	}
}

// Convert all alert fields to annotations.
func (h *Handler) Annotations(f of.ACIFaultRaw) map[string]string {
	// refs:
	// * https://stackoverflow.com/a/18927729
	// * https://play.golang.org/p/_zSICvw562P

	v := reflect.ValueOf(f)

	annotations := make(map[string]string, v.NumField())

	for i := 0; i < v.NumField(); i++ {
		c := strcase.CaseString(v.Type().Field(i).Name)
		snakeCaseOldKey := c.ToSnake()
		key := fmt.Sprintf("fault_%s", snakeCaseOldKey)
		value := v.Field(i).String()
		annotations[key] = value
	}

	return annotations
}

func (h *Handler) GetAlertConfig(fault of.ACIFaultRaw) (string, *of.AlertConfig, error) {
	// refs:
	// * https://stackoverflow.com/a/18927729
	// * https://play.golang.org/p/_zSICvw562P

	// Loop through alerts.yaml:alerts; if an alertConfig mentions the current fault code,
	// return the alertName and alertConfig and break the loop.

	faultCode := ""
	matchFound := false

	v1 := reflect.ValueOf(h.ac)

	for _, e1 := range v1.MapKeys() {
		alertName := e1.Interface().(string)
		newAlertConfig := v1.MapIndex(e1).Interface().(of.AlertConfig)

		// Loop through the nested faults to complete the above goal.
		v2 := reflect.ValueOf(newAlertConfig.Faults)
		for _, e2 := range v2.MapKeys() {
			// ref: https://stackoverflow.com/a/51977415
			faultCode = e2.Interface().(string)
			if faultCode == fault.Code {
				matchFound = true

				// Commenting this out, but we may want to use it in the future.
				newAlertConfigFault := v2.MapIndex(e2).Interface().(of.AlertConfigFault)
				log.Infof("%+v", newAlertConfigFault)
				break
			}
		}

		if matchFound {
			return alertName, &newAlertConfig, nil
		}
	}

	// We didn't find anything...
	err := fmt.Errorf("getAlertConfig() was unable to locate a alert map for fault code: %s", fault.Code)
	return faultCode, &of.AlertConfig{}, err
}

func (h *Handler) Shutdown() error {
	return h.server.Shutdown()
}
