package v2_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	v1 "github.com/cisco-cx/of/pkg/v1"
	mib_registry "github.com/cisco-cx/of/wrap/mib/v1"
)

// Test mibs data.
func mibRegistry(t *testing.T) *mib_registry.MibRegistry {
	mibs := map[string]v1.Mib{
		".1.3.6.1.2.1.1.3.0": v1.Mib{
			Name: "oid1",
		},
		".1.3.6.1.6.3.1.1.4.1.0": v1.Mib{
			Name: "oid2",
		},
		".1.3.6.1.4.1.8164.2.44": v1.Mib{
			Name: "oid3",
		},
		".1.3.6.1.4.1.8164.2.45": v1.Mib{
			Name: "oid4",
		},
		".1.3.6.1.4.1.65000.1.1.1.1.1": v1.Mib{
			Name: "oid5",
		},
	}

	mr := mib_registry.New()
	err := mr.Load(mibs)
	require.NoError(t, err)
	return mr
}

// Test trap vars data.
func trapVars() *[]v1.TrapVar {
	return &[]v1.TrapVar{
		v1.TrapVar{
			Oid:   ".1.3.6.1.2.1.1.3.0",
			Type:  "Timeticks",
			Value: "(123) 0:00:01.23",
		},
		v1.TrapVar{
			Oid:   ".1.3.6.1.6.3.1.1.4.1",
			Type:  "OID",
			Value: ".1.3.6.1.4.1.8164.2.13",
		},
		v1.TrapVar{
			Oid:   ".1.3.6.1.6.3.1.1.4.1.0",
			Type:  "OID",
			Value: ".1.3.6.1.4.1.8164.2.44",
		},
		v1.TrapVar{
			Oid:   ".1.3.6.1.4.1.8164.2.44",
			Type:  "STRING",
			Value: "foo",
		},
		v1.TrapVar{
			Oid:   ".1.3.6.1.6.3.1.1.4.1.1",
			Type:  "OID",
			Value: ".1.3.6.1.4.1.8164.2.45",
		},
		v1.TrapVar{
			Oid:   ".1.3.6.1.4.1.8164.2.45",
			Type:  "OID",
			Value: ".1.3.6.1.4.1.65000.1.1.1.1.1",
		},
		v1.TrapVar{
			Oid:   ".1.3.6.1.4.1.65000.1.1.1.1.1",
			Type:  "STRING",
			Value: "bar",
		},
	}
}
