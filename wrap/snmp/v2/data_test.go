package v2_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	of "github.com/cisco-cx/of/pkg/v2"
	mib_registry "github.com/cisco-cx/of/wrap/mib/v2"
)

// Test mibs data.
func mibRegistry(t *testing.T) *mib_registry.MibRegistry {
	mibs := map[string]of.Mib{
		".1.3.6.1.2.1.1.3.0": of.Mib{
			Name: "oid1",
		},
		".1.3.6.1.6.3.1.1.4.1.0": of.Mib{
			Name: "oid2",
		},
		".1.3.6.1.4.1.8164.2.44": of.Mib{
			Name: "oid3",
		},
		".1.3.6.1.4.1.8164.2.45": of.Mib{
			Name: "oid4",
		},
		".1.3.6.1.4.1.65000.1.1.1.1.1": of.Mib{
			Name: "oid5",
		},
	}

	mr := mib_registry.New()
	err := mr.Load(mibs)
	require.NoError(t, err)
	return mr
}

// Test trap vars data.
func trapVars() *[]of.TrapVar {
	return &[]of.TrapVar{
		of.TrapVar{
			Oid:   ".1.3.6.1.2.1.1.3.0",
			Type:  "Timeticks",
			Value: "(123) 0:00:01.23",
		},
		of.TrapVar{
			Oid:   ".1.3.6.1.6.3.1.1.4.1",
			Type:  "OID",
			Value: ".1.3.6.1.4.1.8164.2.13",
		},
		of.TrapVar{
			Oid:   ".1.3.6.1.6.3.1.1.4.1.0",
			Type:  "OID",
			Value: ".1.3.6.1.4.1.8164.2.44",
		},
		of.TrapVar{
			Oid:   ".1.3.6.1.4.1.8164.2.44",
			Type:  "STRING",
			Value: "foo",
		},
		of.TrapVar{
			Oid:   ".1.3.6.1.6.3.1.1.4.1.1",
			Type:  "OID",
			Value: ".1.3.6.1.4.1.8164.2.45",
		},
		of.TrapVar{
			Oid:   ".1.3.6.1.4.1.8164.2.45",
			Type:  "OID",
			Value: ".1.3.6.1.4.1.65000.1.1.1.1.1",
		},
		of.TrapVar{
			Oid:   ".1.3.6.1.4.1.65000.1.1.1.1.1",
			Type:  "STRING",
			Value: "bar",
		},
	}
}
