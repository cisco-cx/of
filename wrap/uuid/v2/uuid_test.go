package v2_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	of "github.com/cisco-cx/of/pkg/v2"
	uuid "github.com/cisco-cx/of/wrap/uuid/v2"
)

// Enforce interface implementation.
func TestUUIDInterface(t *testing.T) {
	var _ of.UUIDGen = &uuid.UUID{}
}

// Test UUID generator.
func TestUUIDGeneration(t *testing.T) {
	u := uuid.UUID{}
	uuid := u.UUID()
	require.NotEmpty(t, uuid)
	require.Len(t, uuid, 36)
}

// Test Fixed UUID generator.
func TestFixedUUIDGeneration(t *testing.T) {
	u := uuid.FixedUUID{}
	uuid := u.UUID()
	require.Equal(t, uuid, "9dcc77fc-dda5-4edf-a683-64f2589036d6")
}
