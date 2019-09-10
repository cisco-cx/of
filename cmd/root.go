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

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	homedir "github.com/cisco-cx/of/wrap/go-homedir/v1alpha1"
)

var cfgFile string

// rootCmd represents the command that runs when no subcommands are called.
var rootCmd = &cobra.Command{
	Use:   "of",
	Short: "Observability Framework",
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
	// rootCmd.PersistentFlags().String("foo", "", "A help for foo")
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	// Define configuration settings.
	cobra.OnInitialize(initConfig)

	// TODO: Decide if we want a global config file.
	// Define persistent flags that run upon calling any action.
	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.of.yaml)")

	// Define local flags that only run upon calling this action.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
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

	// Read in environment variables that match.
	viper.AutomaticEnv()

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
