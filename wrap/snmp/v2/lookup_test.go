package v2_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	of_snmp "github.com/cisco-cx/of/pkg/v2/snmp"
	logger "github.com/cisco-cx/of/wrap/logrus/v2"
	snmp "github.com/cisco-cx/of/wrap/snmp/v2"
	yaml "github.com/cisco-cx/of/wrap/yaml/v2"
)

// Enforce Lookup interface
func TestLookupInterface(t *testing.T) {
	var _ of_snmp.Lookup = &snmp.Lookup{}
}

// Build lookup map
func TestBuild(t *testing.T) {
	build(t)
}

// Find oid in lookup map
func TestFind(t *testing.T) {
	// Prepare snmp.V2Config
	lookup := build(t)
	vars := trapVars()

	configs, err := lookup.Find(vars)
	require.NoError(t, err)
	require.ElementsMatch(t, []string{"epc", "nso", "device_not_found"}, configs)
}

// Helper to build lookup map
func build(t *testing.T) *snmp.Lookup {
	// Prepare snmp.V2Config
	r := strings.NewReader(YamlContent)
	cfg := yaml.Configs{}
	err := cfg.Decode(r)
	require.NoError(t, err)
	lookup := snmp.Lookup{Configs: of_snmp.V2Config(cfg), MR: mibRegistry(t), Log: logger.New()}
	lookup.Build()
	return &lookup
}
