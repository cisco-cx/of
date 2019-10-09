package v1_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	of "github.com/cisco-cx/of/pkg/v1"
	alertmanager "github.com/cisco-cx/of/wrap/alertmanager/v1"
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
