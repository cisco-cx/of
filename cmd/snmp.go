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
	of_v2 "github.com/cisco-cx/of/pkg/v2"
	mib_registry "github.com/cisco-cx/of/wrap/mib/v2"
	profile "github.com/cisco-cx/of/wrap/profile/v1"
	snmp "github.com/cisco-cx/of/wrap/snmp/v2"
)

// cmdSNMP returns the `snmp` command.
func cmdSNMP() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "snmp",
		Short: "Commands for the SNMP integration",
	}
	// Define flags and configuration settings.
	// cmd.PersistentFlags().String("foo", "", "A help for foo")
	// cmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	return cmd
}

// cmdSNMPHandler returns the `snmp handler` command.
func cmdSNMPHandler() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "handler",
		Short: "Start the SNMP handler",
		Run:   runSNMPHandler,
	}

	// Define flags and configuration settings.
	// cmd.PersistentFlags().String("foo", "", "A help for foo")
	// cmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	cmd.Flags().String("listen-address", "localhost:80", "host:port on which to listen, for SNMP trap events.")
	cmd.Flags().String("am-address", "http://localhost:9093", "AlertManager's URL")
	cmd.Flags().Duration("am-timeout", 1*time.Second, "Alertmanager timeout  (default: 10s)")
	cmd.Flags().String("mibs-dir", "", "Path to MIBs directory.")
	cmd.Flags().String("cache-file", "none", "Path to MIBs cache file.")
	cmd.Flags().String("config-dir", "", "Path to directory containing configs.")
	cmd.Flags().Bool("throttle", true, "Trottle posts to Alertmanager (default: true)")
	cmd.Flags().Int("post-time", 300, "Approx time in ms, that it takes to HTTP POST to AM. (default: 300)")
	cmd.Flags().Int("sleep-time", 100, "Time in ms, to sleep between HTTP POST to AM. (default: 100)")
	cmd.Flags().Int("send-time", 60000, "Time in ms, to complete HTTP POST to AM. (default: 60000)")

	// Enable ENV to set flag values.
	// Ex: ENV AM_URL will set the value for --am-url.
	// Precedence: CLI flag, os.ENV, default value set while defining cmd.Flags().
	viper.BindPFlags(cmd.Flags())
	return cmd
}

// cmdSNMPMIBsProcessor returns the `snmp handler` command.
func cmdSNMPMIBsProcessor() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mib-preprocess",
		Short: "Pre-process JSON MIBs into a single JSON file",
		Run:   runMibsPreProcess,
	}

	// Define flags and configuration settings.
	// cmd.PersistentFlags().String("foo", "", "A help for foo")
	// cmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	cmd.Flags().String("mibs-dir", "", "Path to MIBs directory.")
	cmd.Flags().String("cache-file", "none", "Path to MIBs cache file.")

	// Enable ENV to set flag values.
	// Ex: ENV AM_URL will set the value for --am-url.
	// Precedence: CLI flag, os.ENV, default value set while defining cmd.Flags().
	viper.BindPFlags(cmd.Flags())
	return cmd
}

// Entry point for ./of snmp mib-preprocess.
func runMibsPreProcess(cmd *cobra.Command, args []string) {
	// Start the profiler and defer stopping it until the program exits.
	defer profile.Start().Stop()

	logv2.WithField("info", infoSvc).Infof("snmp mib-preprocess called")

	checkRequiredFlags(cmd)

	SNMPMIBsDir := viper.GetString("mibs-dir")
	cacheFile := viper.GetString("cache-file")

	if SNMPMIBsDir == "" || cacheFile == "none" {
		logv2.Errorf("Please specify a mibs-dir and cache-file.")
		return
	}

	readerMIB := &mib_registry.MIBHandler{
		MapMIB: make(map[string]of_v2.MIB),
	}

	err := readerMIB.LoadJSONFromDir(SNMPMIBsDir)
	if err != nil {
		logv2.WithError(err).Fatalf("Failed to load MIBs from MIBS dir.")
	}
	err = readerMIB.WriteCacheToFile(cacheFile)
	if err != nil {
		logv2.WithError(err).Fatalf("Failed to load MIBs from cache.")
	}
}

// Entry point for ./of snmp handler.
func runSNMPHandler(cmd *cobra.Command, args []string) {
	// Start the profiler and defer stopping it until the program exits.
	defer profile.Start().Stop()

	logv2.WithField("info", infoSvc).Infof("snmp handler called")

	config := SNMPConfig(cmd)
	service, err := snmp.NewService(logv2, config)
	if err != nil {
		logv2.WithError(err).Fatalf("Failed to init SNMP service.")
	}

	handler := &snmp.Handler{
		Config: config,
		SNMP:   service,
		Log:    logv2,
	}

	handler.Run()
}

// Returns  &of.SNMPConfig{} based on CLI flags and ENV.
func SNMPConfig(cmd *cobra.Command) *of_v2.SNMPConfig {
	checkRequiredFlags(cmd)
	cfg := &of_v2.SNMPConfig{}
	cfg.ListenAddress = viper.GetString("listen-address")
	cfg.AMAddress = viper.GetString("am-address")
	cfg.AMTimeout = viper.GetDuration("am-timeout")
	cfg.SNMPMibsDir = viper.GetString("mibs-dir")
	cfg.CacheFile = viper.GetString("cache-file")
	cfg.ConfigDir = viper.GetString("config-dir")
	cfg.Version = infoSvc.String()

	cfg.Throttle = viper.GetBool("throttle")
	cfg.PostTime = viper.GetInt("post-time")
	cfg.SleepTime = viper.GetInt("sleep-time")
	cfg.SendTime = viper.GetInt("send-time")

	if strings.HasPrefix(cfg.AMAddress, "http") == false {
		logv2.Fatalf("AM URL must begin with http/https")
	}

	return cfg
}

// init creates and adds the `snmp` command with its subcommands.
func init() {
	// Create the `snmp` command.
	cmd := cmdSNMP()

	// Create and add the subcommands of the `snmp` command.
	cmd.AddCommand(cmdSNMPHandler())
	cmd.AddCommand(cmdSNMPMIBsProcessor())

	// Add the `snmp` command to the root command.
	rootCmd.AddCommand(cmd)
}
