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
	"strings"
	"time"

	"github.com/mitchellh/mapstructure"
	of "github.com/cisco-cx/of/lib/v1alpha1"
	acigo "github.com/cisco-cx/of/wrap/acigo/v1alpha1"
	alertmanager "github.com/cisco-cx/of/wrap/alertmanager/v1alpha1"
	http "github.com/cisco-cx/of/wrap/http/v1alpha1"
	logger "github.com/cisco-cx/of/wrap/logrus/v1alpha1"
	prometheus "github.com/cisco-cx/of/wrap/prometheus/client_golang/v1alpha1"
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
	Log      *logger.Logger
}

func (h *Handler) Run() {

	h.InitHandler()
	srv := of.Server{
		Addr:         h.Config.ListenAddress,
		ReadTimeout:  h.Config.ACITimeout,
		WriteTimeout: h.Config.ACITimeout,
	}

	go func() {
		for {
			h.PushAlerts()
			time.Sleep(time.Duration(h.Config.CycleInterval) * time.Second)
		}
	}()

	h.server = http.NewServer(srv)

	h.server.HandleFunc("/", func(w of.ResponseWriter, r of.Request) {
		fmt.Fprint(w, h.Config.Version)
	})

	h.server.Handle("/metrics", prometheus.NewHandler())
	err := h.server.ListenAndServe()
	if err != nil {
		h.Log.WithError(err).Fatalf("Failed to listen at %s", h.Config.ListenAddress)
	}
}

func (h *Handler) InitHandler() {

	h.ac = &yaml.Alerts{}
	h.LoadConfig(h.ac, h.Config.AlertsCFGFile)
	h.sc = &yaml.Secrets{}
	h.LoadConfig(h.sc, h.Config.SecretsCFGFile)

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
			h.Log.WithError(err).Fatalf("Failed to init counter, %s", name)
		}
	}
}

// Pull ACI faults and forward to Alertmanager.
func (h *Handler) PushAlerts() {

	var alerts []*alertmanager.Alert
	var err error
	h.counters[notificationCycleCount].Incr()
	h.Log.Debugf("Running APIC -> AlertManager notification cycle. (cycle-sleep-seconds=%d)\n", h.Config.CycleInterval)

	h.counters[apicConnectAttemptCount].Incr()
	faults, err := h.Aci.Faults()
	if err != nil {
		h.counters[apicConnectErrorCount].Incr()
		h.Log.WithError(err).Errorf("Failed to get faults.")
		return
	}

	alerts, err = h.FaultsToAlerts(faults)
	if err != nil {
		h.Log.WithError(err).Errorf("Failed to convert faults to alerts.")
		return
	}

	// Notify AlertManager. if we have any alerts
	if len(alerts) > 0 {
		h.counters[amConnectAttemptCount].Incr()
		err = h.Ams.Notify(alerts)
		if err != nil {
			h.counters[amConnectErrorCount].Incr()
			h.Log.Errorf("Notification cycle failed. Will retry in %d, %s\n", h.Config.CycleInterval, err.Error())
			return
		}
	} else {
		h.Log.Errorf("No faults found")
	}

	h.Log.Debugf("Notification cycle succeeded. Sleeping for %d seconds.\n", h.Config.CycleInterval)
}

// Convert acigo faults to Alertmanager alerts.
func (h *Handler) FaultsToAlerts(faults []of.Map) ([]*alertmanager.Alert, error) {
	var alerts []*alertmanager.Alert
	for _, mapFault := range faults {

		// Decode fault into struct.
		f := of.ACIFaultRaw{}
		mapstructure.Decode(mapFault, &f)
		fp := acigo.FaultParser{f, h.Log}
		h.counters[faultsScrapedCount].Incr()

		// If this is in alerts.yaml:dropped_faults, skip it.
		if _, drop := h.ac.DroppedFaults[strings.ToUpper(fp.Fault.Code)]; drop {
			h.Log.Debugf("Dropping fault: %s\n", f)
			h.counters[faultsDroppedCount].Incr()
			continue
		}

		// Create alert boilerplate.
		alert := alertmanager.NewAlert(f)

		// Get an integer representation of fault severity for numerical comparison.
		s, err := acigo.NewACIFaultSeverityRaw(fp.Fault.Severity)
		if err != nil {
			h.Log.Errorf("Failed to parse severity, %s", err.Error())
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
			h.Log.Debugf("Found matching fault code in alertsConfig.Alerts.")
			alert.Labels["alertname"] = of.LabelValue(alertName)
			alert.Labels["alert_severity"] = of.LabelValue(newAlertConfig.AlertSeverity)
		} else {
			h.counters[faultsUnmatchedCount].Incr()
			h.Log.Debugf("%s\n", err)
			h.Log.Debugf("Setting default alertname and severity for fault code: %s\n", f.Code)

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
			h.Log.Errorf("Failed to parse severity, %s", err.Error())
			return nil, err
		}

		faultSeverityThresholdLevel := s.ID()
		if faultSeverityLevel < faultSeverityThresholdLevel {
			alert.EndsAt = faultLastTransition
		}

		alert.Labels[amAlertFingerprintLabel] = of.LabelValue(alert.Fingerprint())

		// Debug sample code.
		b := []byte{}
		b, err = json.Marshal(alert)
		if err != nil {
			return nil, err
		}
		h.Log.WithField("Alert name", alert.Labels["alertname"]).
			WithField("Fingerprint", alert.Labels[amAlertFingerprintLabel]).
			Infof("Alert generated.")
		h.Log.Debugf("Alert generated: %s\n", b)

		alerts = append(alerts, alert)
		h.counters[alertsGeneratedCount].Incr()
	}
	return alerts, nil

}

// Wrapper to read a file into an implementation of of.Decoder.
func (h *Handler) LoadConfig(cfg of.Decoder, fileName string) {

	f, err := os.Open(fileName)
	if err != nil {
		h.Log.WithError(err).Fatalf("Failed to open file. config-file : %s", fileName)
	}

	err = cfg.Decode(f)
	if err != nil {
		h.Log.WithError(err).Fatalf("Failed to decode alerts config file.")
	}
}

func (h *Handler) GetAlertConfig(fault of.ACIFaultRaw) (string, *of.AlertConfig, error) {
	// Loop through alerts.yaml:alerts; if an alertConfig mentions the current fault code,
	// return the alertName and alertConfig and break the loop.

	for alertName, alertConfig := range h.ac.Alerts {
		for code, _ := range alertConfig.Faults {
			if code == fault.Code {
				return alertName, &alertConfig, nil
			}
		}
	}

	// We didn't find anything...
	err := fmt.Errorf("getAlertConfig() was unable to locate a alert map for fault code: %s", fault.Code)
	return "", &of.AlertConfig{}, err
}

func (h *Handler) Shutdown() error {
	return h.server.Shutdown()
}
