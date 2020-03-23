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
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/cisco-cx/of/info"
	of "github.com/cisco-cx/of/pkg/v1"
	aci "github.com/cisco-cx/of/wrap/aci/v1"
	acigo "github.com/cisco-cx/of/wrap/acigo/v1"
	alertmanager "github.com/cisco-cx/of/wrap/alertmanager/v1"
	net "github.com/cisco-cx/of/wrap/net/v1"
	profile "github.com/cisco-cx/of/wrap/profile/v1"
)

const staticLabelUsage = "Static labels to be added with each alert posted to Alertmanager. Expected format : 'label1=value1,label2=value2'"

// cmdACI returns the `aci` command.
func cmdACI() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "aci",
		Short: "Commands for the ACI integration",
	}
	return cmd
}

// cmdACIHandler returns the `aci handler` command.
func cmdACIHandler() *cobra.Command {
	cmd := &cobra.Command{
		Use:                "handler",
		Short:              "Start the ACI handler",
		Run:                runACIHandler,
		DisableFlagParsing: true,
	}

	return cmd
}

// Entry point for ./of aci handler.
func runACIHandler(cmd *cobra.Command, args []string) {
	// Start the profiler and defer stopping it until the program exits.
	defer profile.Start().Stop()

	log.WithField("info", infoSvc).Infof("aci handler called")

	cmd.Flags().String("aci-listen-address", "localhost:9011", "host:port on which to listen, for metrics scraping")
	cmd.Flags().Int("aci-cycle-interval", 60, "Number of seconds to sleep between APIC -> AM notification cycles (default: 60)")
	cmd.Flags().String("aci-am-url", "", "[Required] AlertManager's URL")
	cmd.Flags().SetAnnotation("aci-am-url", "required", []string{"true"})
	cmd.Flags().String("aci-host", "", "[Required] ACI host (Value is ignored when --aci-enable-consul is set")
	cmd.Flags().SetAnnotation("aci-host", "required", []string{"true"})
	cmd.Flags().String("aci-user", "", "[Required] ACI username")
	cmd.Flags().SetAnnotation("aci-user", "required", []string{"true"})
	cmd.Flags().String("aci-password", "", "[Required] ACI password")
	cmd.Flags().SetAnnotation("aci-password", "required", []string{"true"})
	cmd.Flags().String("aci-alerts-config", "alerts.yaml", "Alerts config file (default: alerts.yaml)")
	cmd.Flags().String("aci-secrets-config", "secrets.yaml", "Secrets config file (default: secrets.yaml)")
	cmd.Flags().Duration("aci-timeout", 10*time.Second, "ACI Read/Write timeout  (default: 10s)")
	cmd.Flags().String("aci-static-labels", "None", staticLabelUsage)
	cmd.Flags().Bool("aci-throttle", true, "Trottle posts to Alertmanager (default: true)")
	cmd.Flags().Int("aci-post-time", 300, "Approx time in ms, that it takes to HTTP POST to AM. (default: 300)")
	cmd.Flags().Int("aci-sleep-time", 100, "Time in ms, to sleep between HTTP POST to AM. (default: 100)")
	cmd.Flags().Int("aci-send-time", 60000, "Time in ms, to complete HTTP POST to AM. (default: 60000)")
	cmd.Flags().Bool("aci-enable-consul", false, "Whether to use consul for host discovery (default: false)")

	checkRequiredFlags(cmd, args, []string{})

	config := ACIConfig(cmd)
	handler := &aci.Handler{Config: config, Log: log}
	handler.Aci = &acigo.ACIService{ACIConfig: config, Logger: log}
	handler.Ams = &alertmanager.AlertService{AmURL: config.AmURL, Version: config.Version}
	handler.Run()
}

// Returns  &ACI{} based on CLI flags and ENV.
func ACIConfig(cmd *cobra.Command) *of.ACIConfig {
	cfg := &of.ACIConfig{}
	cfg.Application = info.Program
	cfg.ListenAddress = viper.GetString("aci-listen-address")
	cfg.CycleInterval = viper.GetInt("aci-cycle-interval")
	cfg.AmURL = viper.GetString("aci-am-url")
	cfg.ACIHost = viper.GetString("aci-host")

	cfg.AlertsCFGFile = viper.GetString("aci-alerts-config")
	cfg.SecretsCFGFile = viper.GetString("aci-secrets-config")

	cfg.User = viper.GetString("aci-user")
	cfg.Pass = viper.GetString("aci-password")
	cfg.Version = infoSvc.String()
	cfg.SourceHostname, cfg.SourceAddress = VerifiedHost(cfg.ACIHost)

	cfg.ACITimeout = viper.GetDuration("aci-timeout")

	cfg.Throttle = viper.GetBool("aci-throttle")
	cfg.PostTime = viper.GetInt("aci-post-time")
	cfg.SleepTime = viper.GetInt("aci-sleep-time")
	cfg.SendTime = viper.GetInt("aci-send-time")

	cfg.ConsulEnabled = viper.GetBool("aci-enable-consul")

	if strings.HasPrefix(cfg.AmURL, "http") == false {
		log.Fatalf("aci-am-url must begin with http/https")
	}

	staticLabels := viper.GetString("aci-static-labels")
	if staticLabels != "" && staticLabels != "None" {
		m := make(of.LabelMap)
		labelItems := strings.Split(staticLabels, ",")
		for _, labelItem := range labelItems {
			kvs := strings.Split(labelItem, "=")
			if len(kvs) != 2 {
				log.Fatalf("%s, given : %s", staticLabelUsage, labelItem)
			}
			m[of.LabelName(kvs[0])] = of.LabelValue(kvs[1])
		}
		cfg.StaticLabels = m
	}

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

// init creates and adds the `aci` command with its subcommands.
func init() {
	// Create the `aci` command.
	cmd := cmdACI()

	// Create and add the subcommands of the `aci` command.
	cmd.AddCommand(cmdACIHandler())

	// Add the `aci` command to the root command.
	rootCmd.AddCommand(cmd)
}
