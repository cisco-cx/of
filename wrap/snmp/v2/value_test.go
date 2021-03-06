package v2_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	of_snmp "github.com/cisco-cx/of/pkg/v2/snmp"
	snmp "github.com/cisco-cx/of/wrap/snmp/v2"
)

// Enforce Lookup interface
func TestValInterface(t *testing.T) {
	var _ of_snmp.ValueGenerator = &snmp.Value{}
}

// Testing Value
func TestValue(t *testing.T) {
	v := newValue(t)
	val, err := v.Value(".1.3.6.1.2.1.1.3.0")
	require.NoError(t, err)
	require.Equal(t, "(123) 0:00:01.23", val)
}

// Testing Value with given As type.
func TestValueAs(t *testing.T) {
	types := map[of_snmp.As]string{
		of_snmp.Value:            ".1.3.6.1.4.1.8164.2.45",
		of_snmp.ValueStr:         ".1.3.6.1.4.1.8164.2.oid4",
		of_snmp.ValueStrShort:    "oid4",
		of_snmp.OidValue:         ".1.3.6.1.4.1.65000.1.1.1.1.1",
		of_snmp.OidValueStr:      ".1.3.6.1.4.1.65000.1.1.1.1.oid5",
		of_snmp.OidValueStrShort: "oid5",
	}
	value := newValue(t)
	for k, v := range types {
		val, err := value.ValueAs(".1.3.6.1.6.3.1.1.4.1.1", k)
		require.NoError(t, err)
		require.Equal(t, v, val)
	}
}

// Testing Value string with numerical OID
func TestValueStr(t *testing.T) {
	v := newValue(t)
	val, err := v.ValueStr(".1.3.6.1.6.3.1.1.4.1.0")
	require.NoError(t, err)
	require.Equal(t, ".1.3.6.1.4.1.8164.2.oid3", val)
}

// Testing Value string with none numerical OID
func TestValueStrFail(t *testing.T) {
	v := newValue(t)
	_, err := v.ValueStr(".1.3.6.1.2.1.1.3.0")
	require.Error(t, err)
}

// Testing Value short string
func TestValueStrShort(t *testing.T) {
	v := newValue(t)
	val, err := v.ValueStrShort(".1.3.6.1.6.3.1.1.4.1.0")
	require.NoError(t, err)
	require.Equal(t, "oid3", val)
}

// Testing OID Value
func TestOIDValue(t *testing.T) {
	v := newValue(t)
	val, err := v.OIDValue(".1.3.6.1.4.1.8164.2.45")
	require.NoError(t, err)
	require.Equal(t, "bar", val)
}

// Testing OID missing in traps
func TestOIDValueFail(t *testing.T) {
	v := newValue(t)
	_, err := v.OIDValue(".1.3.6.1.4.1.65000.1.1.1.1.1")
	require.Error(t, err)
}

// Testing OID Value string
func TestOIDValueStr(t *testing.T) {
	v := newValue(t)
	val, err := v.OIDValueStr(".1.3.6.1.6.3.1.1.4.1.1")
	require.NoError(t, err)
	require.Equal(t, ".1.3.6.1.4.1.65000.1.1.1.1.oid5", val)
}

// Testing OID Value short string
func TestOIDValueStrShort(t *testing.T) {
	v := newValue(t)
	val, err := v.OIDValueStrShort(".1.3.6.1.6.3.1.1.4.1.1")
	require.NoError(t, err)
	require.Equal(t, "oid5", val)
}

// Initialize snmp.Value
func newValue(t *testing.T) *snmp.Value {
	return snmp.NewValue(trapVars(), mibRegistry(t))
}
