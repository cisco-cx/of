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
	"github.com/spf13/cobra"
)

// cmdVersion returns the `version` command.
func cmdVersion() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Display version information",
		Run:   runVersion,
	}
	return cmd
}

func runVersion(cmd *cobra.Command, args []string) {
	log.WithField("info", infoSvc).Infof("version called")
}

// init creates and adds the `version` command
func init() {
	// Create the `version` command.
	cmd := cmdVersion()

	// Add the `version` command to the root command.
	rootCmd.AddCommand(cmd)
}
