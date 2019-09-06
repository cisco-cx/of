package v1alpha1_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/cisco-cx/of/lib/v1alpha1"
	configv1alpha1 "github.com/cisco-cx/of/wrap/config/v1alpha1"
)

func TestAlertsLoader(t *testing.T) {

	r := strings.NewReader(`
apic:
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

	cfg := configv1alpha1.Alerts{}
	cfg.Load(r)
	require.EqualValues(t, cfg, expected)
}
