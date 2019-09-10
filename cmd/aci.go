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

	"github.com/spf13/cobra"
)

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

// cmdACIServer returns the `aci server` command.
func cmdACIServer() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "handler",
		Short: "Start the ACI handler",
		Run: runACIServer,
	}
	// Define flags and configuration settings.
	// cmd.PersistentFlags().String("foo", "", "A help for foo")
	// cmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	return cmd
}

func runACI(cmd *cobra.Command, args []string) {
	fmt.Println("aci called")
}

func runACIServer(cmd *cobra.Command, args []string) {
	fmt.Println("aci handler called")
}

// init creates and adds the `aci` command with its subcommands.
func init() {
	// Create the `aci` command.
	cmd := cmdACI()

	// Create and add the subcommands of the `aci` command.
	cmd.AddCommand(cmdACIServer())

	// Add the `aci` command to the root command.
	rootCmd.AddCommand(cmd)
}
