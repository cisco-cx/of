// Copyright 2019 Cisco Systems, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package v1_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	of "github.com/cisco-cx/of/pkg/v1"
	yaml "github.com/cisco-cx/of/wrap/yaml/v1"
)

// Enforce interface implementation.
func TestAlertsInterface(t *testing.T) {
	var _ of.Decoder = &yaml.Alerts{}
	var _ of.Encoder = &yaml.Alerts{}
}

// Ensure yaml decodes Alerts
func TestAlertsDecoder(t *testing.T) {

	r := strings.NewReader(
		`apic:
  alert_severity_threshold: minor
defaults:
  alert_severity: error
alerts:
  apicFabricSelectorIssuesConfig:  
    alert_severity: error
    faults:  
      F0020:  
        fault_name: fltFabricSelectorIssuesConfig-failed  
dropped_faults:  
  F3104:
    fault_name: xxx  
  F2100:
    fault_name: yyy  
  F675299:
    fault_name: fsmFailHcloudHealthUpdateSyncHealth`)

	expected := of.Alerts{
		APIC: of.AlertsConfigAPIC{
			AlertSeverityThreshold: "minor",
		},
		Defaults: of.AlertsConfigDefaults{
			AlertSeverity: "error",
		},
		DroppedFaults: map[string]of.AlertsConfigDroppedFault{
			"F3104": of.AlertsConfigDroppedFault{
				FaultName: "xxx",
			},
			"F2100": of.AlertsConfigDroppedFault{
				FaultName: "yyy",
			},
			"F675299": of.AlertsConfigDroppedFault{
				FaultName: "fsmFailHcloudHealthUpdateSyncHealth",
			},
		},
		Alerts: map[string]of.AlertConfig{
			"apicFabricSelectorIssuesConfig": {
				AlertSeverity: "error",
				Faults: map[string]of.AlertConfigFault{
					"F0020": of.AlertConfigFault{
						FaultName: "fltFabricSelectorIssuesConfig-failed",
					},
				},
			},
		},
	}

	cfg := yaml.Alerts{}
	cfg.Decode(r)
	require.EqualValues(t, expected, cfg)
}

// Ensure yaml encodes Alerts
func TestAlertsEncoder(t *testing.T) {

	expected := strings.Trim(`
apic:
  alert_severity_threshold: minor
defaults:
  alert_severity: error
dropped_faults:
  F2100:
    fault_name: yyy
  F3104:
    fault_name: xxx
  F675299:
    fault_name: fsmFailHcloudHealthUpdateSyncHealth
alerts:
  apicFabricSelectorIssuesConfig:
    alert_severity: error
    faults:
      F0020:
        fault_name: fltFabricSelectorIssuesConfig-failed`, "\n")

	cfg := yaml.Alerts{
		APIC: of.AlertsConfigAPIC{
			AlertSeverityThreshold: "minor",
		},
		Defaults: of.AlertsConfigDefaults{
			AlertSeverity: "error",
		},
		DroppedFaults: map[string]of.AlertsConfigDroppedFault{
			"F3104": of.AlertsConfigDroppedFault{
				FaultName: "xxx",
			},
			"F2100": of.AlertsConfigDroppedFault{
				FaultName: "yyy",
			},
			"F675299": of.AlertsConfigDroppedFault{
				FaultName: "fsmFailHcloudHealthUpdateSyncHealth",
			},
		},
		Alerts: map[string]of.AlertConfig{
			"apicFabricSelectorIssuesConfig": {
				AlertSeverity: "error",
				Faults: map[string]of.AlertConfigFault{
					"F0020": of.AlertConfigFault{
						FaultName: "fltFabricSelectorIssuesConfig-failed",
					},
				},
			},
		},
	}

	buf := bytes.NewBuffer(nil)
	cfg.Encode(buf)
	require.EqualValues(t, expected, strings.Trim(string(buf.Bytes()), "\n"))
}
