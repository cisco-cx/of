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
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/cisco-cx/of/info"
	homedir "github.com/cisco-cx/of/wrap/go-homedir/v1"
	informer "github.com/cisco-cx/of/wrap/informer/v1"
	logger "github.com/cisco-cx/of/wrap/logrus/v1"
	loggerv2 "github.com/cisco-cx/of/wrap/logrus/v2"
)

var cfgFile string
var log = logger.New()
var logv2 = loggerv2.New()

// Start a shared info service.
var infoSvc = informer.NewInfoService(
	info.Program,
	info.License,
	info.URL,
	info.BuildUser,
	info.BuildDate,
	info.Language,
	info.LanguageVersion,
	info.Version,
	info.Revision,
	info.Branch,
)

// rootCmd represents the command that runs when no subcommands are called.
var rootCmd = &cobra.Command{
	Use:   "of",
	Short: "Observability Framework",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {

		viper.BindPFlags(cmd.Flags())
		err := cmd.Flags().Parse(args)

		// skip "unknown flag" error messages for CLI parsing since each subcommand can add it's own flags
		if err != nil && !strings.Contains(err.Error(), "unknown flag") {
			return err
		}
		logLevel := viper.GetString("log-level")
		jsonLogging := viper.GetBool("json-logging")
		log.SetLevel(logLevel)
		logv2.SetLevel(logLevel)
		if jsonLogging == true {
			log.EnableJSONLogging()
			logv2.EnableJSONLogging()
		}
		log.Infof("Logging Enabled. Level : %s", log.LogLevel())
		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}

// "init is called after all the variable declarations in the package have
//  evaluated their initializers, and those are evaluated only after all the
//  imported packages have been initialized."
//
// source: https://golang.org/doc/effective_go.html#init
func init() {
	// Define flags and configuration settings.
	rootCmd.PersistentFlags().String("log-level", "info", "Log Level")
	rootCmd.PersistentFlags().Bool("json-logging", false, "Enable logging in JSON format. (default: false)")

	viper.BindPFlags(rootCmd.PersistentFlags())
	// Define configuration settings.
	cobra.OnInitialize(initConfig)

}

// initConfig reads in config file and ENV variables if set.
//
// TODO: If we have config files per subcommand instead of one global config,
// let's rework this.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".of" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".of")
	}

	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	// Read in environment variables that match.
	viper.AutomaticEnv()

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

// Check if all flags are set and print the values for each flag.
func checkRequiredFlags(cmd *cobra.Command, args []string, optArgs []string) {
	// Enable ENV to set flag values.
	// Ex: ENV AM_URL will set the value for --am-url.
	// Precedence: CLI flag, os.ENV, default value set while defining cmd.Flags().
	viper.BindPFlags(cmd.Flags())
	parseErr := cmd.Flags().Parse(args)

	optMap := make(map[string]bool)
	optMap["help"] = true
	for _, v := range optArgs {
		optMap[v] = true
	}

	fmt.Println("=============== Current Setting ===============")
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		log.Debugf("flag : %s\n", f.Name)
		if _, ok := optMap[f.Name]; ok == false {
			if isFlagSet(f) == false {
				cmd.Usage()
				fmt.Printf("Required : %s\n", f.Name)
				if parseErr != nil {
					fmt.Println(parseErr.Error())
				}
				os.Exit(1)
			}
		}
		fmt.Printf("%16s : %-24v // %s\n", f.Name, viper.Get(f.Name), f.Usage)
	})
	fmt.Println("===============================================")
}

// Check if flag is set.
func isFlagSet(f *pflag.Flag) bool {
	// If flag is "required" and is still no value is passed from commang line, return false
	requiredAnnotation := f.Annotations["required"]
	if len(requiredAnnotation) == 0 {
		return true
	} else if requiredAnnotation[0] == "true" && viper.Get(f.Name) == f.DefValue {
		// f.Changed gets set only when passed on command line, thus environment variables, even if set correctly, fail the check
		return false
	}
	return true
}
