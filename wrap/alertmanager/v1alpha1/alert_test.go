package v1alpha1_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	of "github.com/cisco-cx/of/lib/v1alpha1"
	alertmanager "github.com/cisco-cx/of/wrap/alertmanager/v1alpha1"
)

// Enforce interface implementation.
func TestAlertInterface(t *testing.T) {
	var _ of.Fingerprinter = &alertmanager.Alert{}
}

// Test finger printing.
func TestAlertFingerPrint(t *testing.T) {
	a := alertmanager.NewAlert(of.ACIFaultRaw{})
	require.Equal(t, "cbf29ce484222325", a.Fingerprint())
}
