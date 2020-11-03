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

package v1

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	of "github.com/cisco-cx/of/pkg/v1"
	aci_config "github.com/cisco-cx/of/pkg/v1/aci"
	of_v2 "github.com/cisco-cx/of/pkg/v2"
	acigo "github.com/cisco-cx/of/wrap/acigo/v1"
	alertmanager "github.com/cisco-cx/of/wrap/alertmanager/v1"
	http "github.com/cisco-cx/of/wrap/http/v2"
	logger "github.com/cisco-cx/of/wrap/logrus/v1"
	net "github.com/cisco-cx/of/wrap/net/v1"
	prometheus "github.com/cisco-cx/of/wrap/prometheus/client_golang/v2"
	yaml "github.com/cisco-cx/of/wrap/yaml/v1"
	consul "github.com/hashicorp/consul/api"
	"github.com/mitchellh/mapstructure"
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
	apicLoginAttemptCount   = "apic_login_attempt_total"
	apicLoginErrorCount     = "apic_login_error_count"
	alertsGeneratedCount    = "alerts_generated_count"
	faultsDroppedCount      = "faults_dropped_count"
	faultsScrapedCount      = "faults_scraped_count"
	faultsMatchedCount      = "faults_matched_count"
	faultsUnmatchedCount    = "faults_unmatched_count"
	faultsUnknownIgnored    = "faults_unknown_ignored_count"
	notificationCycleCount  = "notification_cycle_count"
	nodeEnriched            = "nodes_enriched_count"
)

type Handler struct {
	Config      *of.ACIConfig
	counters    map[string]*prometheus.Counter
	counterVecs map[string]*prometheus.CounterVec
	server      *http.Server
	Aci         *acigo.ACIClient
	Ams         *alertmanager.AlertService
	ac          *yaml.Alerts
	sc          *yaml.Secrets
	Log         *logger.Logger
}

func (h *Handler) Run() {

	h.InitHandler()

	go func() {
		nodes := h.Config.ACIHosts
		for {
			if h.Config.ConsulEnabled {
				nodes = h.GetConsulNodes()
				h.Log.Infof("ACI nodes received from consul: %s\n", nodes)
			}
			if len(nodes) <= 0 {
				h.Log.Warningf("Empty host list.\n")
			}
			startTime := time.Now()
			for i := 0; i < len(nodes); i++ {
				h.Log.Debugf("Fetching errors from %s\n", nodes[i])
				h.Aci.SetHost(nodes[i])
				h.PushAlerts()
			}
			elapsedTime := time.Duration(time.Since(startTime).Seconds())
			sleepTime := time.Duration(h.Config.CycleInterval) - elapsedTime
			if sleepTime > 0 {
				h.Log.Debugf("Sleeping for %d seconds.\n", sleepTime)
				time.Sleep(sleepTime * time.Second)
			}
		}
	}()

	httpConfig := of_v2.HTTPConfig{
		ListenAddress: h.Config.ListenAddress,
		ReadTimeout:   h.Config.ACITimeout,
		WriteTimeout:  h.Config.ACITimeout,
	}
	h.server = http.NewServer(&httpConfig, h.Config.Application)

	h.server.HandleFunc("/", func(w of_v2.ResponseWriter, r of_v2.Request) {
		fmt.Fprint(w, h.Config.Version)
	})

	err := h.server.ListenAndServe()
	if err != nil {
		h.Log.WithError(err).Fatalf("Failed to listen at %s", h.Config.ListenAddress)
	}

	<-make(chan bool)
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
		apicLoginAttemptCount: &prometheus.Counter{Namespace: h.Config.Application, Name: apicLoginAttemptCount,
			Help: "Number of times we tried to login to APIC."},
		apicLoginErrorCount: &prometheus.Counter{Namespace: h.Config.Application, Name: apicLoginErrorCount,
			Help: "Number of errors encountered while trying to login to APIC."},
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
		faultsUnknownIgnored: &prometheus.Counter{Namespace: h.Config.Application, Name: faultsUnknownIgnored,
			Help: "Number of times we could not find an alertConfig that mentioned the encountered fault code and the fault was ignored."},
		notificationCycleCount: &prometheus.Counter{Namespace: h.Config.Application, Name: notificationCycleCount,
			Help: "Number of times we tried ran the notification cycle loop."},
	}

	for name, c := range h.counters {
		err := c.Create()
		if err != nil {
			h.Log.WithError(err).Fatalf("Failed to init counter, %s", name)
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
				Namespace: h.Config.Application,
				Name:      nodeEnriched,
				Help:      "Number of times we enriched topology.",
			},
			labels: []string{"fault_dn", "node_dn", "enriched", "node", "role"},
		},
	}

	h.counterVecs = make(map[string]*prometheus.CounterVec)
	for _, vi := range vis {
		err := vi.vector.Create(vi.labels)
		if err != nil {
			h.Log.WithError(err).Fatalf("Failed to init counterVec, %s", vi.vector.Name)
		}
		h.counterVecs[vi.vector.Name] = vi.vector
	}
}

// GetConsulNodes Lists the nodes from consul, matching given service and node metadata
func (h *Handler) GetConsulNodes() []string {
	config := consul.Config{Address: h.sc.Consul.Host}
	consulClient, err := consul.NewClient(&config)
	if err != nil {
		h.Log.WithError(err).Fatalf("Failed to create Consul client")
	}

	queryOptions := consul.QueryOptions{NodeMeta: h.sc.Consul.NodeMeta}
	if queryOptions.NodeMeta == nil {
		queryOptions.NodeMeta = make(map[string]string)
	}
	if strings.TrimSpace(h.Config.ConsulACIGroupHost) != "" {
		queryOptions.NodeMeta["aci-group-host"] = h.Config.ConsulACIGroupHost
	}
	service, _, err := consulClient.Catalog().Service(h.sc.Consul.Service, "", &queryOptions)

	if err != nil {
		h.Log.WithError(err).Fatalf("Failed to get Nodes from the Consul Service")
	}

	nodes := make([]string, 0)
	for _, node := range service {
		nodes = append(nodes, node.Node)
	}

	return nodes
}

// Pull ACI faults and forward to Alertmanager.
func (h *Handler) PushAlerts() {

	var alerts []*alertmanager.Alert
	var err error
	h.counters[notificationCycleCount].Incr()
	h.Log.Debugf("Running APIC -> AlertManager notification cycle.\n")

	h.counters[apicLoginAttemptCount].Incr()
	err = h.Aci.Login()
	if err != nil {
		h.counters[apicLoginErrorCount].Incr()
		h.Log.WithError(err).Errorf("Failed to login.")
		return
	}
	defer h.Aci.Logout()

	nodes, err := h.Aci.NodeList()
	if err != nil {
		h.counters[apicConnectErrorCount].Incr()
		h.Log.WithError(err).Errorf("Failed to get nodes.")
		return
	}
	h.counters[apicConnectAttemptCount].Incr()
	faults, err := h.Aci.Faults()
	h.Log.Debugf("Retrieved %d faults from ACI", len(faults))
	if err != nil {
		h.counters[apicConnectErrorCount].Incr()
		h.Log.WithError(err).Errorf("Failed to get faults.")
		return
	}

	alerts, err = h.FaultsToAlerts(faults, nodes)
	h.Log.Debugf("Converted faults to %d alerts", len(alerts))
	if err != nil {
		h.Log.WithError(err).Errorf("Failed to convert faults to alerts.")
		return
	}

	// Notify AlertManager. if we have any alerts
	if len(alerts) > 0 {
		// Send alerts[start:end] to Alertmanager.
		sendFunc := func(start int, end int) {
			h.Log.Debugf("Sending alerts from %d to %d\n", start, end)
			h.counters[amConnectAttemptCount].Incr()
			err = h.Ams.Notify(alerts[start:end])
			if err != nil {
				h.counters[amConnectErrorCount].Incr()
				h.Log.Errorf("Notification cycle failed. Will retry in %d, %s\n", h.Config.CycleInterval, err.Error())
				return
			}
		}

		h.Throttle(len(alerts), sendFunc)
	} else {
		h.Log.Errorf("No faults found")
	}
}

// Divide alerts into smaller chunks and spread posting to Alertmanager over h.Config.SendTime milliseconds.
func (h *Handler) Throttle(totalCount int, f func(int, int)) {

	// Send all alerts in a single post to Alertmanager, if Throttle is disabled
	// or h.Config.SendTime is less than time needed for a single post.
	if h.Config.Throttle == false || h.Config.SendTime <= h.Config.PostTime+h.Config.SleepTime {
		f(0, totalCount)
		h.Log.Infof("Throttle disabled, sending all alerts")
		return
	}

	// Max num. of requests that can be send in h.Config.SendTime.
	maxRequests := h.Config.SendTime / (h.Config.PostTime + h.Config.SleepTime)
	start := 0
	if totalCount > maxRequests {
		chunkSize := totalCount / maxRequests

		end := chunkSize
		for end <= totalCount {
			f(start, end)
			start = end
			end = start + chunkSize
			h.Log.Infof("Throttling alerts, sleeping for %d seconds", h.Config.SleepTime)
			time.Sleep(time.Duration(h.Config.SleepTime) * time.Millisecond)
		}
	}

	// Handle condition where totalCount is not divisible by maxRequests.
	if start < totalCount {
		f(start, totalCount)
		h.Log.Infof("Throttling alerts, sleeping for %d seconds", h.Config.SleepTime)
		time.Sleep(time.Duration(h.Config.SleepTime) * time.Millisecond)
	}
}

// Convert acigo faults to Alertmanager alerts.
func (h *Handler) FaultsToAlerts(faults []of.Map, nodes map[string]map[string]interface{}) ([]*alertmanager.Alert, error) {
	var alerts []*alertmanager.Alert
	for _, mapFault := range faults {
		h.Log.Tracef("Processing fault: %s\n", mapFault)

		// Decode fault into struct.
		f := of.ACIFaultRaw{}
		err := mapstructure.Decode(mapFault, &f)
		if err != nil {
			h.Log.Errorf("Failed to decode map structure with error: %s, structure: %+v", err.Error(), mapFault)
			return nil, err
		}
		h.Log.Tracef("Decoded fault: %+v\n", f)

		fp := acigo.FaultParser{Fault: f, Log: h.Log}
		h.counters[faultsScrapedCount].Incr()

		// If this is in alerts.yaml:dropped_faults, skip it.
		if _, drop := h.ac.DroppedFaults[strings.ToUpper(fp.Fault.Code)]; drop {
			h.Log.Debugf("Dropping fault: %s\n", f)
			h.counters[faultsDroppedCount].Incr()
			continue
		}

		// Create alert boilerplate.
		alert := alertmanager.NewAlert(f)
		h.Log.Tracef("New Alert: %+v", alert)

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
			h.Log.Errorf("Failed to parse created date, %s", err.Error())
			return nil, err
		}
		faultLastTransition, err := fp.LastTransition()
		if err != nil {
			h.Log.Errorf("Failed to parse last transition date, %s", err.Error())
			return nil, err
		}

		// Set key values on the alert.
		alert.StartsAt = faultCreated
		alert.GeneratorURL = fmt.Sprintf(apicFaultHelpURL, f.Code)
		alert.Labels["cluster_name"] = of.LabelValue(h.sc.APIC.Cluster.Name)

		// Adding custom labels.
		for l, v := range h.Config.StaticLabels {
			alert.Labels[l] = v
		}

		sub_id, _ := fp.SubID()
		alert.Labels["sub_id"] = of.LabelValue(sub_id)
		alert.Annotations["source_hostname"], alert.Annotations["source_address"] = h.VerifiedHost(h.Aci.GetHost())

		// Try to find this fault code in alerts config.
		alertName, newAlertConfig, err := h.GetAlertConfig(f)
		if err == nil && alertName != "" {
			h.counters[faultsMatchedCount].Incr()
			h.Log.Debugf("Found matching fault code in alertsConfig.Alerts.")
			alert.Labels["alertname"] = of.LabelValue(alertName)
			alert.Labels["alert_severity"] = of.LabelValue(newAlertConfig.AlertSeverity)
		} else if h.ac.APIC.DropUnknownAlerts {
			// the alert wasn't found and we are ignoring unknown alerts
			h.Log.Debugf("Ignoring unknown fault code=%s, rule=%s\n", f.Code, f.Rule)
			h.counters[faultsUnknownIgnored].Incr()
			continue
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

		// Enrich topology
		if h.ac.EnrichTopology.Enabled == true {
			h.EnrichTopology(alert, nodes)
		}

		alert.Labels[amAlertFingerprintLabel] = of.LabelValue(alert.Fingerprint())

		// Debug sample code.
		b := []byte{}
		b, err = json.Marshal(alert)
		if err != nil {
			h.Log.Errorf("Failed to marshal to JSON, %s", err.Error())
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

// Do a forward and reverse lookup to verify the ACI Host.
// If DNS entry is found, Hostname and IP from DNS Query is returned
// else aciHost is returned
func (h *Handler) VerifiedHost(aciHost string) (string, string) {

	hostname := aciHost
	ipAddr := aciHost

	// DNS reverse lookup
	ip, err := net.NewIP(aciHost)
	if err != nil {
		h.Log.WithError(err).Errorf("")
	}

	hostnames, err := ip.Hostnames()
	if err != nil {
		h.Log.WithError(err).Errorf("Failed to find hostname.")
	}
	if len(hostnames) == 0 {
		h.Log.Errorf("No reverse lookup available for %s", ip.String())
	} else {
		hostname = string(hostnames[0])
	}

	// DNS forward lookup
	host, err := net.NewHostname(hostname)
	if err != nil {
		h.Log.WithError(err).Fatalf("")
	}

	var ips []of.IP
	ips, err = host.IPv6()
	if err != nil || len(ips) == 0 {
		ips, err = host.IPv4()
	}

	if err == nil && len(ips) != 0 {
		ipAddr = string(ips[len(ips)-1])
	}

	return hostname, ipAddr
}

// Wrapper to update Labels with node and subsystem.
func (h *Handler) EnrichTopology(alert *alertmanager.Alert, nodes map[string]map[string]interface{}) {
	faultDN := alert.Annotations["fault_dn"]
	if h.ac.EnrichTopology.DefaultSubsystem != "" {
		alert.Labels["subsystem"] = of.LabelValue(h.ac.EnrichTopology.DefaultSubsystem)
	}

	if strings.HasPrefix(faultDN, "topology") {
		nodeDN := strings.Join(strings.Split(faultDN, "/")[0:3], "/")
		if nodeInfo, ok := nodes[nodeDN]; ok == true {
			h.counterVecs[nodeEnriched].Incr(map[string]string{
				"fault_dn": faultDN,
				"node_dn":  nodeDN,
				"enriched": "true",
				"node":     nodeInfo["name"].(string),
				"role":     nodeInfo["role"].(string),
			})
			if h.ac.EnrichTopology.NodeLabel != "" {
				nodeName := nodeInfo["name"].(string)
				if strings.HasSuffix(nodeName, h.ac.EnrichTopology.NodeTLD) == false {
					nodeName = fmt.Sprintf("%s%s", nodeName, h.ac.EnrichTopology.NodeTLD)
				}
				alert.Labels[of.LabelName(h.ac.EnrichTopology.NodeLabel)] = of.LabelValue(nodeName)
			}
			if role, ok := h.ac.EnrichTopology.Subsystems[nodeInfo["role"].(string)]; ok == true {
				alert.Labels["subsystem"] = of.LabelValue(role)
			}
		} else {
			h.counterVecs[nodeEnriched].Incr(map[string]string{
				"fault_dn": faultDN,
				"node_dn":  nodeDN,
				"enriched": "false",
				"node":     "unknown",
				"role":     "unknown",
			})
			h.Log.WithFields(map[string]interface{}{
				"fault_dn": faultDN,
				"node_dn":  nodeDN,
			}).Errorf("Failed to match fault with node.")
		}
	}

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

func (h *Handler) GetAlertConfig(fault of.ACIFaultRaw) (string, *aci_config.AlertConfig, error) {
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
	return "", &aci_config.AlertConfig{}, err
}

func (h *Handler) Shutdown() error {
	return h.server.Shutdown()
}
