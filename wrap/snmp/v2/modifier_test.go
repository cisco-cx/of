package v2_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	of_snmp "github.com/cisco-cx/of/pkg/v2/snmp"
	snmp "github.com/cisco-cx/of/wrap/snmp/v2"
)

// Enforce implementing of_snmp.Modifier.
func TestModifierInterface(t *testing.T) {
	var _ of_snmp.Modifier = &snmp.Modifier{}
}

// Test Mod Apply.
func TestModApply(t *testing.T) {
	//m := snmp.NewModifier(newValueModifier(t))

	mods := []of_snmp.Mod{
		of_snmp.Mod{
			Key:   "Vendor",
			Value: "Cisco",
			Type:  of_snmp.Set,
		},
		of_snmp.Mod{
			Key:   "System",
			Value: "nso",
			Type:  of_snmp.Set,
		},

		of_snmp.Mod{
			Oid:   ".1.3.6.1.6.3.1.1.4.1.1",
			Type:  of_snmp.Copy,
			ToKey: "Value",
			As:    of_snmp.Value,
		},
		of_snmp.Mod{
			Oid:   ".1.3.6.1.6.3.1.1.4.1.1",
			Type:  of_snmp.Copy,
			ToKey: "ValueStr",
			As:    of_snmp.ValueStr,
		},
		of_snmp.Mod{
			Oid:   ".1.3.6.1.6.3.1.1.4.1.1",
			Type:  of_snmp.Copy,
			ToKey: "ValueStrShort",
			As:    of_snmp.ValueStrShort,
		},
		of_snmp.Mod{
			Oid:   ".1.3.6.1.6.3.1.1.4.1.1",
			Type:  of_snmp.Copy,
			ToKey: "OidValue",
			As:    of_snmp.OidValue,
		},
		of_snmp.Mod{
			Oid:   ".1.3.6.1.6.3.1.1.4.1.1",
			Type:  of_snmp.Copy,
			ToKey: "OidValueStr",
			As:    of_snmp.OidValueStr,
		},
		of_snmp.Mod{
			Oid:   ".1.3.6.1.6.3.1.1.4.1.1",
			Type:  of_snmp.Copy,
			ToKey: "OidValueStrShort",
			As:    of_snmp.OidValueStrShort,
		},
	}

	expectedLabel := map[string]string{
		"Vendor":           "Cisco",
		"System":           "nso",
		"Value":            ".1.3.6.1.4.1.8164.2.45",
		"ValueStr":         ".1.3.6.1.4.1.8164.2.oid4",
		"ValueStrShort":    "oid4",
		"OidValue":         ".1.3.6.1.4.1.65000.1.1.1.1.1",
		"OidValueStr":      ".1.3.6.1.4.1.65000.1.1.1.1.oid5",
		"OidValueStrShort": "oid5",
	}

	label := make(map[string]string)
	m := snmp.Modifier{Map: &label, V: newValueModifier(t)}
	err := m.Apply(mods)
	require.NoError(t, err)
	require.Equal(t, expectedLabel, label)
}

// Test Mod Copy with map and value present in map.
func TestModCopyWithMap(t *testing.T) {

	mod := of_snmp.Mod{
		Oid:   ".1.3.6.1.6.3.1.1.4.1.1",
		Type:  of_snmp.Copy,
		ToKey: "OidValueStrShort",
		As:    of_snmp.OidValueStrShort,
		Map: map[string]string{
			"oid5": "Found in map",
		},
	}

	expectedLabel := map[string]string{
		"OidValueStrShort": "Found in map",
	}

	label := make(map[string]string)
	m := snmp.Modifier{Map: &label, V: newValueModifier(t)}
	err := m.Copy(mod)
	require.NoError(t, err)
	require.Equal(t, expectedLabel, label)

}

// Test Mod Copy with map and value not present.
func TestModCopyWithMapEntryMissing(t *testing.T) {

	mod := of_snmp.Mod{
		Oid:   ".1.3.6.1.6.3.1.1.4.1.1",
		Type:  of_snmp.Copy,
		ToKey: "OidValueStrShort",
		As:    of_snmp.OidValueStrShort,
		Map: map[string]string{
			"some key": "Won't find oid5 here",
		},
	}

	expectedLabel := map[string]string{}

	label := make(map[string]string)
	m := snmp.Modifier{Map: &label, V: newValueModifier(t)}
	err := m.Copy(mod)
	require.NoError(t, err)
	require.Equal(t, expectedLabel, label)
}

// Test Mod Copy with of_snmp.Send as onError
func TestModWithOnErrorSend(t *testing.T) {

	mods := []of_snmp.Mod{
		of_snmp.Mod{
			Key:   "Vendor",
			Value: "Cisco",
			Type:  of_snmp.Set,
		},

		// This should fail but, set should still happen.
		of_snmp.Mod{
			Oid:     "OidToError",
			Type:    of_snmp.Copy,
			ToKey:   "OidValueStrShort",
			As:      of_snmp.OidValueStrShort,
			OnError: of_snmp.Send, // This is the default behaviour if not set.
		},
	}

	expectedLabel := map[string]string{
		"Vendor": "Cisco",
	}

	label := make(map[string]string)
	m := snmp.Modifier{Map: &label, V: newValueModifier(t)}
	err := m.Apply(mods)
	require.NoError(t, err)
	require.Equal(t, expectedLabel, label)
}

// Test Mod Copy with of_snmp.Send as onError
func TestModWithOnErrorDrop(t *testing.T) {

	mods := []of_snmp.Mod{
		of_snmp.Mod{
			Key:   "Vendor",
			Value: "Cisco",
			Type:  of_snmp.Set,
		},

		// This should fail but, set should still happen.
		of_snmp.Mod{
			Oid:     "OidToError",
			Type:    of_snmp.Copy,
			ToKey:   "OidValueStrShort",
			As:      of_snmp.OidValueStrShort,
			OnError: of_snmp.Drop,
		},
	}

	label := make(map[string]string)
	m := snmp.Modifier{Map: &label, V: newValueModifier(t)}
	err := m.Apply(mods)
	require.Error(t, err)
}

// Test Mod Set.
func TestModSet(t *testing.T) {

	expectedLabel := map[string]string{
		"vendor": "cisco",
		"system": "epc",
	}

	label := make(map[string]string)
	label["vendor"] = "cisco"

	mod := of_snmp.Mod{
		Type:  of_snmp.Set,
		Key:   "system",
		Value: "epc",
	}

	m := snmp.Modifier{Map: &label, V: newValueModifier(t)}
	err := m.Set(mod)
	require.NoError(t, err)
	require.Equal(t, expectedLabel, label)
}

// Test Mod with wrong operation.
func TestModWrongOp(t *testing.T) {

	label := make(map[string]string)

	mod := of_snmp.Mod{
		Type: of_snmp.Copy,
	}

	m := snmp.Modifier{Map: &label, V: newValueModifier(t)}
	err := m.Set(mod)
	require.Error(t, err)
}

// Test Mod with missing key.
func TestModMissingKey(t *testing.T) {

	label := make(map[string]string)

	mod := of_snmp.Mod{
		Type: of_snmp.Copy,
	}

	m := snmp.Modifier{Map: &label, V: newValueModifier(t)}
	err := m.Copy(mod)
	require.Error(t, err)
}

// Initialize snmp.Value.
func newValueModifier(t *testing.T) *snmp.Value {
	return snmp.NewValue(trapVars(), mibRegistry(t))
}
