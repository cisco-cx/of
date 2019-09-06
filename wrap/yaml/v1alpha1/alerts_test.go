// Copyright 2019 Cisco Systems, Inc.
//
// This work incorporates works covered by the following notice:
//
// Copyright 2018 The Prometheus Authors
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

package v1alpha1_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/cisco-cx/of/lib/v1alpha1"
	yaml "github.com/cisco-cx/of/wrap/yaml/v1alpha1"
)

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

	expected := v1alpha1.Alerts{
		APIC: v1alpha1.AlertsConfigAPIC{
			AlertSeverityThreshold: "minor",
		},
		Defaults: v1alpha1.AlertsConfigDefaults{
			AlertSeverity: "error",
		},
		DroppedFaults: map[string]v1alpha1.AlertsConfigDroppedFault{
			"F3104": v1alpha1.AlertsConfigDroppedFault{
				FaultName: "xxx",
			},
			"F2100": v1alpha1.AlertsConfigDroppedFault{
				FaultName: "yyy",
			},
			"F675299": v1alpha1.AlertsConfigDroppedFault{
				FaultName: "fsmFailHcloudHealthUpdateSyncHealth",
			},
		},
		Alerts: map[string]v1alpha1.AlertConfig{
			"apicFabricSelectorIssuesConfig": {
				AlertSeverity: "error",
				Faults: map[string]v1alpha1.AlertConfigFault{
					"F0020": v1alpha1.AlertConfigFault{
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
		APIC: v1alpha1.AlertsConfigAPIC{
			AlertSeverityThreshold: "minor",
		},
		Defaults: v1alpha1.AlertsConfigDefaults{
			AlertSeverity: "error",
		},
		DroppedFaults: map[string]v1alpha1.AlertsConfigDroppedFault{
			"F3104": v1alpha1.AlertsConfigDroppedFault{
				FaultName: "xxx",
			},
			"F2100": v1alpha1.AlertsConfigDroppedFault{
				FaultName: "yyy",
			},
			"F675299": v1alpha1.AlertsConfigDroppedFault{
				FaultName: "fsmFailHcloudHealthUpdateSyncHealth",
			},
		},
		Alerts: map[string]v1alpha1.AlertConfig{
			"apicFabricSelectorIssuesConfig": {
				AlertSeverity: "error",
				Faults: map[string]v1alpha1.AlertConfigFault{
					"F0020": v1alpha1.AlertConfigFault{
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
