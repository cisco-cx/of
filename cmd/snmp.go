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
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	of_v2 "github.com/cisco-cx/of/pkg/v2"
	mib_registry "github.com/cisco-cx/of/wrap/mib/v2"
	profile "github.com/cisco-cx/of/wrap/profile/v1"
	prometheus "github.com/cisco-cx/of/wrap/prometheus/client_golang/v2"
	snmp "github.com/cisco-cx/of/wrap/snmp/v2"
	watcher "github.com/cisco-cx/of/wrap/watcher/v2"
)

// cmdSNMP returns the `snmp` command.
func cmdSNMP() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "snmp",
		Short: "Commands for the SNMP integration",
	}
	return cmd
}

// cmdSNMPHandler returns the `snmp handler` command.
func cmdSNMPHandler() *cobra.Command {
	cmd := &cobra.Command{
		Use:                "handler",
		Short:              "Start the SNMP handler",
		Run:                runSNMPHandler,
		DisableFlagParsing: true,
	}

	return cmd
}

// cmdSNMPMIBsProcessor returns the `snmp handler` command.
func cmdSNMPMIBsProcessor() *cobra.Command {
	cmd := &cobra.Command{
		Use:                "mib-preprocess",
		Short:              "Pre-process JSON MIBs into a single JSON file",
		Run:                RunMibsPreProcess,
		DisableFlagParsing: true,
	}
	return cmd
}

// Entry point for ./of snmp mib-preprocess.
func RunMibsPreProcess(cmd *cobra.Command, args []string) {
	// Start the profiler and defer stopping it until the program exits.
	defer profile.Start().Stop()

	logv2.WithField("info", infoSvc).Infof("snmp mib-preprocess called")

	// Define flags and configuration settings.
	cmd.Flags().String("mibs-dir", "", "Path to MIBs directory.")
	cmd.Flags().String("cache-file", "none", "Path to MIBs cache file.")

	checkRequiredFlags(cmd, args, []string{})

	SNMPMIBsDir := viper.GetString("mibs-dir")
	cacheFile := viper.GetString("cache-file")

	if SNMPMIBsDir == "" || cacheFile == "none" {
		logv2.Errorf("Please specify a input-dir and output-file.")
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
	ParseSNMPHandlerFlags(cmd, args)
	config := SNMPConfig(cmd)
	logv2.Infof("Starting SNMP service")

	fs, err := watcher.NewPath(config.ConfigDir, logv2)
	if err != nil {
		logv2.WithError(err).Fatalf("Failed to init watcher for config dir.")
	}

	cntr, cntrVec := snmp.InitCounters(config.Application, logv2)
	handler := initSNMPHandler(config, cntr, cntrVec)
	cntrVec[snmp.HandlerRestarted].Incr(map[string]string{
		"op_type": "start",
	})
	go handler.Run()

	var restartTimer *time.Timer

	restartFunc := func() {
		logv2.Infof("Shutting down SNMP handler")
		cntrVec[snmp.HandlerRestarted].Incr(map[string]string{
			"op_type": "shutdown",
		})
		handler.Shutdown()
		handler = initSNMPHandler(config, cntr, cntrVec)
		logv2.Infof("Starting SNMP handler")
		cntrVec[snmp.HandlerRestarted].Incr(map[string]string{
			"op_type": "start",
		})
		handler.Run()
		restartTimer = nil
	}

	err = fs.Watch()
	if err != nil {
		logv2.WithError(err).Fatalf("Failed to watch config dir.")
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	for {
		select {
		case _ = <-fs.Changed:
			// Might get multiple events, since we are watching a directory.
			// Delaying and restarting the service only once.
			logv2.Infof("Config change detected.")
			if restartTimer == nil {
				logv2.Debugf("Started timer to restart handler.")
				restartTimer = time.AfterFunc(2*time.Second, restartFunc)
			} else {
				logv2.Debugf("Delaying handler restart.")
				if !restartTimer.Stop() {
					<-restartTimer.C
				}
				restartTimer.Reset(time.Second)
			}
		case _ = <-signalChan:
			logv2.Infof("Process killed, shutting down.")
			handler.Shutdown()
			os.Exit(0)
		}
	}

}

func initSNMPHandler(config *of_v2.SNMPConfig, cntr map[string]*prometheus.Counter, cntrVec map[string]*prometheus.CounterVec) *snmp.Handler {
	service, err := snmp.NewService(logv2, config, cntr, cntrVec)
	if err != nil {
		logv2.WithError(err).Fatalf("Failed to init SNMP service.")
	}

	handler := &snmp.Handler{
		Config: config,
		SNMP:   service,
		Log:    logv2,
	}
	return handler
}

func ParseSNMPHandlerFlags(cmd *cobra.Command, args []string) {

	// Define flags and configuration settings.
	cmd.Flags().String("listen-address", "localhost:80", "host:port on which to listen, for SNMP trap events.")
	cmd.Flags().String("am-address", "http://localhost:9093", "AlertManager's URL")
	cmd.Flags().Duration("am-timeout", 1*time.Second, "Alertmanager timeout  (default: 10s)")
	cmd.Flags().String("mibs-dir", "none", "Path to MIBs directory.")
	cmd.Flags().String("cache-file", "none", "Path to MIBs cache file.")
	cmd.Flags().String("config-dir", "", "Path to directory containing configs.")
	cmd.Flags().Bool("throttle", true, "Trottle posts to Alertmanager (default: true)")
	cmd.Flags().Int("post-time", 300, "Approx time in ms, that it takes to HTTP POST to AM. (default: 300)")
	cmd.Flags().Int("sleep-time", 100, "Time in ms, to sleep between HTTP POST to AM. (default: 100)")
	cmd.Flags().Int("send-time", 10000, "Time in ms, to complete HTTP POST to AM. (default: 10000)")
	cmd.Flags().Bool("dry-run", false, "Log generated alerts, instead of sending to Alertmanager. (default: false)")
	cmd.Flags().Bool("log-unknown", false, "Log unknown alerts at info level. (default: false)")
	cmd.Flags().Bool("forward-unknown", false, "send unknown alerts to Alertmanager. (default: false)")
	checkRequiredFlags(cmd, args, []string{"dry-run", "log-unknown", "forward-unknown"})
}

// Returns  &of.SNMPConfig{} based on CLI flags and ENV.
func SNMPConfig(cmd *cobra.Command) *of_v2.SNMPConfig {
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
	cfg.DryRun = viper.GetBool("dry-run")
	cfg.LogUnknown = viper.GetBool("log-unknown")
	cfg.ForwardUnknown = viper.GetBool("forward-unknown")

	if strings.HasPrefix(cfg.AMAddress, "http") == false {
		logv2.Fatalf("AM URL must begin with http/https")
	}

	if cfg.SNMPMibsDir == "none" && cfg.CacheFile == "none" {
		logv2.Fatalf("Please specify a mibs-dir or cache-file.")
	}

	// Setting namespace for counters.
	cfg.Application = "of_snmp_handler"
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
