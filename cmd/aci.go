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

package cmd

import (
	"fmt"
	"os"
	"reflect"
	"runtime"
	"time"

	"github.com/prometheus/common/model"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	of "github.com/cisco-cx/of/lib/v1alpha1"
	aci "github.com/cisco-cx/of/wrap/aci/v1alpha1"
	acigo "github.com/cisco-cx/of/wrap/acigo/v1alpha1"
	alertmanager "github.com/cisco-cx/of/wrap/alertmanager/v1alpha1"
	http "github.com/cisco-cx/of/wrap/http/v1alpha1"
	net "github.com/cisco-cx/of/wrap/net/v1alpha1"
	prometheus "github.com/cisco-cx/of/wrap/prometheus/client_golang/v1alpha1"
	yaml "github.com/cisco-cx/of/wrap/yaml/v1alpha1"
)

// Represents ACI settings.
type ACI struct {
	Application    string
	ListenAddress  string
	CycleInterval  int
	AmURL          string
	ACIHost        string
	AlertsCFGFile  string
	SecretsCFGFile string
	user           string
	pass           string
	ac             *yaml.Alerts
	sc             *yaml.Secrets
	sourceHostname string
	sourceAddress  string
	counters       map[string]*prometheus.Counter
	server         *http.Server
}

// Represents Alertmanager alert.
type amAlert struct {
	Labels       model.LabelSet    `json:"labels"`
	Annotations  map[string]string `json:"annotations"`
	StartsAt     time.Time         `json:"startsAt"`
	EndsAt       time.Time         `json:"endsAt"`
	GeneratorURL string            `json:"generatorURL"`
}

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

// Alertmanager alert specific constants.
const (
	apicFaultHelpURL        = "https://pubhub.devnetcloud.com/media/apic-mim-ref-411/docs/FAULT-%s.html"
	amAlertFingerprintLabel = "alert_fingerprint"
)

// App info.
var application = "amapicclient"
var revision = "unset"

// cmdACI returns the `aci` command.
func cmdACI() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "aci",
		Short: "Commands for the ACI integration",
	}
	// Define flags and configuration settings.
	// cmd.PersistentFlags().String("foo", "", "A help for foo")
	// cmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	return cmd
}

// cmdACIHandler returns the `aci handler` command.
func cmdACIHandler() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "handler",
		Short: "Start the ACI handler",
		Run:   runACIHandler,
	}

	// Define flags and configuration settings.
	// cmd.PersistentFlags().String("foo", "", "A help for foo")
	// cmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	cmd.Flags().String("listen-address", "localhost:9011", "host:port on which to listen, for metrics scraping")
	cmd.Flags().Int("cycle-interval", 60, "Number of seconds to sleep between APIC -> AM notification cycles (default: 60)")
	cmd.Flags().String("am-url", "", "AlertManager's URL")
	cmd.Flags().String("aci-host", "", "ACI host")
	cmd.Flags().String("aci-user", "", "ACI username")
	cmd.Flags().String("aci-password", "", "ACI password")
	cmd.Flags().String("alerts-config", "alerts.yaml", "Alerts config file (default: alerts.yaml)")
	cmd.Flags().String("secrets-config", "secrets.yaml", "Secrets config file (default: secrets.yaml)")

	// Enable ENV to set flag values.
	// Ex: ENV AM_URL will set the value for --am-url.
	// Precedence: CLI flag, os.ENV, default value set while defining cmd.Flags().
	viper.BindPFlags(cmd.Flags())
	return cmd
}

// Entry point for ./of aci handler.
func runACIHandler(cmd *cobra.Command, args []string) {
	config := ACIConfig(cmd)
	handler := *&aci.Handler{Config: config}
	handler.Aci = &acigo.ACIService{config}
	handler.Ams = &alertmanager.AlertService{config}
	handler.Run()
}

// Returns  &ACI{} based on CLI flags and ENV.
func ACIConfig(cmd *cobra.Command) *of.ACIConfig {
	checkRequiredFlags(cmd)
	cfg := &of.ACIConfig{}
	cfg.Application = application
	cfg.ListenAddress = viper.GetString("listen-address")
	cfg.CycleInterval = viper.GetInt("cycle-interval")
	cfg.AmURL = viper.GetString("am-url")
	cfg.ACIHost = viper.GetString("aci-host")

	cfg.AlertsCFGFile = viper.GetString("alerts-config")
	cfg.SecretsCFGFile = viper.GetString("secrets-config")

	cfg.User = viper.GetString("aci-user")
	cfg.Pass = viper.GetString("aci-password")
	cfg.Version = fmt.Sprintf("%s %s (%s)", application, revision, runtime.Version())
	cfg.SourceHostname, cfg.SourceAddress = VerifiedHost(cfg.ACIHost)

	t := time.Now()
	zone, offset := t.Zone()
	log.Infof("This machine's timezone and offset are: %s, %d\n", zone, offset)

	return cfg
}

// Do a forward and reverse lookup to verify the ACI Host.
// If DNS entry is found, Hostname and IP from DNS Query is returned
// else aciHost is returned
func VerifiedHost(aciHost string) (string, string) {

	hostname := aciHost
	ipAddr := aciHost

	// DNS reverse lookup
	ip, err := net.NewIP(aciHost)
	if err != nil {
		log.WithError(err).Errorf("")
	}

	hostnames, err := ip.Hostnames()
	if err != nil {
		log.WithError(err).Errorf("Failed to find hostname.")
	}
	if len(hostnames) == 0 {
		log.Errorf("No reverse lookup available for %s", ip.String())
	} else {
		hostname = string(hostnames[0])
	}

	// DNS forward lookup
	host, err := net.NewHostname(hostname)
	if err != nil {
		log.WithError(err).Fatalf("")
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

// Check if flag is set.
func isFlagSet(f *pflag.Flag) bool {
	// If flag value is not set return false
	val := viper.Get(f.Name)
	if val == reflect.Zero(reflect.TypeOf(val)).Interface() {
		return false
	}

	return true
}

// Check if all flags are set and print the values for each flag.
func checkRequiredFlags(cmd *cobra.Command) {
	fmt.Println("=============== Current Setting ===============")
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		log.Debugf("flag : %s\n", f.Name)
		if f.Name != "help" {
			if isFlagSet(f) == false {
				cmd.Usage()
				fmt.Printf("Required : %s\n", f.Name)
				os.Exit(1)
			}
			fmt.Printf("%16s : %-24v // %s\n", f.Name, viper.Get(f.Name), f.Usage)
		}
	})
	fmt.Println("===============================================")
}

// init creates and adds the `aci` command with its subcommands.
func init() {
	// Create the `aci` command.
	cmd := cmdACI()

	// Create and add the subcommands of the `aci` command.
	cmd.AddCommand(cmdACIHandler())

	// Add the `aci` command to the root command.
	rootCmd.AddCommand(cmd)
}
