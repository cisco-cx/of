package v2_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	v1 "github.com/cisco-cx/of/pkg/v1"
	of_snmp "github.com/cisco-cx/of/pkg/v2/snmp"
	mib_registry "github.com/cisco-cx/of/wrap/mib/v1"
	snmp "github.com/cisco-cx/of/wrap/snmp/v2"
	yaml "github.com/cisco-cx/of/wrap/yaml/v2"
)

// Sample configs.
var yamlContent = `epc:
  defaults:
    enabled: true
    source_type: host
    generator_url_prefix: http://www.oid-info.com/get/
    label_mods:
    - type: set
      key: vendor
      value: cisco
    - type: set
      key: subsystem
      value: epc
    - type: copy
      oid: .1.3.6.1.4.1.8164.1.2.1.1.1
      as: value
      to_key: star_slot_num
      on_error: drop
    annotation_mods:
    - type: copy
      oid: .1.3.6.1.6.1.1.1.4.1
      as: value
      to_key: event_oid
    - type: copy
      oid: .1.3.6.1.6.1.1.1.4.1
      as: oid.value-str-short
      to_key: event_name
  alerts:
  - name: starCard
    enabled: true
    label_mods:
    - type: set
      key: alert_severity
      value: error
    firing:
      select:
      - type: equals
        oid: .1.3.6.1.6.3.1.1.4.1
        as: value
        values:
        - .1.3.6.1.4.1.8164.2.13
        - .1.3.6.1.4.1.8164.2.4
        - .1.3.6.1.4.1.8164.2.7
        - .1.3.6.1.4.1.8164.2.44
    clearing:
      select:
      - type: equals
        oid: .1.3.6.1.6.3.1.1.4.1
        as: value
        values:
        - .1.3.6.1.4.1.8164.2.5
        - .1.3.6.1.4.1.8164.2.55
  - name: starCardBootFailed
    label_mods:
    - type: set
      key: alert_severity
      value: error
    firing:
      select:
      - type: equals
        oid: .1.3.6.1.6.3.1.1.4.1
        as: value
        values:
        - .1.3.6.1.4.1.8164.2.9
    clearing:
      select:
      - type: equals
        oid: .1.3.6.1.6.3.1.1.4.1
        as: value
        values:
        - .1.3.6.1.4.1.8164.2.5
        - .1.3.6.1.4.1.8164.2.55
  - name: starCardActive
    label_mods:
    - type: set
      key: alert_severity
      value: informational
    firing:
      select:
      - type: equals
        oid: .1.3.6.1.6.3.1.1.4.1
        as: value
        values:
        - .1.3.6.1.4.1.8164.2.55
nso:
  defaults:
    source_type: cluster   # (host|cluster)...
      # if cluster, you must define defaults.clusters.
    clusters:
      nso1.example.org:  # cluster_name
        source_addresses:  # TODO: enhance this when necessary.
        - 192.168.1.28
        - dead:beef::1
    generator_url_prefix: http://www.oid-info.com/get/
      # numerical OID is appended automatically
    label_mods:
      # Allows promotion from snmpTrapOID information to labels.
      # You cannot promote from annotations to labels.
    - type: set
      key: vendor
      value: cisco
    - type: set
      key: subsystem
      value: nso
    - type: copy
      oid: .1.3.6.1.4.1.24961.2.103.1.1.5.1.2  # tfAlarmType
      as: value
      to_key: alertname
      map:  # is non null, so we're looking up in a map
        alarm-type: nsoAlarmType
        ncs-cluster-alarm: nsoNcsClusterAlarm
        cluster-subscriber-failure: nsoClusterSubcriberFailure
        ncs-dev-manager-alarm: nsoNcsDevManagerAlarm
        ned-live-tree-connection-failure: nsoNedLiveTreeConnectionFailure
        dev-manager-internal-error: nsoDevManagerInternalError
        final-commit-error: nsoFinalCommitError
        commit-through-queue-blocked: nsoCommitThroughQueueBlocked
        abort-error: nsoAbortError
        revision-error: nsoRevisionError
        missing-transaction-id: nsoMissingTransactionId
        configuration-error: nsoConfigurationError
        commit-through-queue-failed: nsoCommitThroughQueueFailed
        connection-failure: nsoConnectionFailure
        out-of-sync: nsoOutOfSync
        ncs-snmp-notification-receiver-alarm: nsoNcsSnmpNotificationReceiverAlarm
        receiver-configuration-error: nsoReceiverConfigurationError
        ncs-package-alarm: nsoNcsPackageAlarm
        package-load-failure: nsoPackageLoadFailure
        package-operation-failure: nsoPackageOperationFailure
        ncs-service-manager-alarm: nsoNcsServiceManagerAlarm
        service-activation-failure: nsoServiceActivationFailure
    - type: copy
      oid: .1.3.6.1.4.1.24961.2.103.1.1.5.1.2  # tfAlarmType
      as: value
      to_key: alert_severity
      map:  # is non null, so we're looking up in a map
        alarm-type: major
        ncs-cluster-alarm: minor
        cluster-subscriber-failure: critical
        ncs-dev-manager-alarm: critical
        ned-live-tree-connection-failure: critical
        dev-manager-internal-error: critical
        final-commit-error: critical
        commit-through-queue-blocked: critical
        abort-error: critical
        revision-error: critical
        missing-transaction-id: critical
        configuration-error: critical
        commit-through-queue-failed: critical
        connection-failure: critical
        out-of-sync: critical
        ncs-snmp-notification-receiver-alarm: critical
        receiver-configuration-error: critical
        ncs-package-alarm: critical
        package-load-failure: critical
        package-operation-failure: critical
        ncs-service-manager-alarm: critical
        service-activation-failure: critical
    annotation_mods: []
    # The service automatically sets annotations.event_type
    # For firing events, annotations.event_type='firing'
    # For clearing events, annotations.event_type='clearing'
    # to_key: event_type
  alerts:
  - name: null  # Auto-set by default.label_mods, need not define
    label_mods:
      # allow promotion from snmpTrapOID information to labels
      # You cannot promote from annotations to labels.
    - type: set
      key: alert_severity
      value: error
    firing:
      select:
      - type: equals
        oid: .1.3.6.1.6.3.1.1.4.1 # snmpTrapOID
        as: value
        values:
        - .1.3.6.1.4.1.24961.2.103.2.0.3  # tfAlarmMinor
        - .1.3.6.1.4.1.24961.2.103.2.0.4  # tfAlarmMajor
        - .1.3.6.1.4.1.24961.2.103.2.0.5  # tfAlarmCritical
      annotation_mods: []  # this is allowed
    clearing:
      select:
      - type: equals
        oid: .1.3.6.1.6.3.1.1.4.1 # snmpTrapOID
        as: value
        values:
        - .1.3.6.1.4.1.24961.2.103.2.0.1  # tfAlarmIndeterminate
        - .1.3.6.1.4.1.24961.2.103.2.0.2  # tfAlarmWarning
        - .1.3.6.1.4.1.24961.2.103.2.0.6  # tfAlarmClear
      annotation_mods: []  # this is allowed`

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
	configs, err := lookup.Find(".1.3.6.1.6.3.1.1.4.1")
	require.NoError(t, err)
	require.Equal(t, []string{"epc"}, configs)
}

// Helper to build lookup map
func build(t *testing.T) *snmp.Lookup {
	// Prepare snmp.V2Config
	r := strings.NewReader(yamlContent)
	cfg := yaml.Configs{}
	err := cfg.Decode(r)
	require.NoError(t, err)
	lookup := snmp.Lookup{Configs: of_snmp.V2Config(cfg), V: newValueLookup(t)}
	lookup.Build()
	return &lookup
}

// Initialize snmp.Value
func newValueLookup(t *testing.T) *snmp.Value {
	return snmp.NewValue(trapVars(), mibRegistry(t))
}

// Test mibs data.
func mibRegistry(t *testing.T) *mib_registry.MibRegistry {
	mibs := map[string]v1.Mib{
		".1.3.6.1.2.1.1.3.0": v1.Mib{
			Name: "oid1",
		},
		".1.3.6.1.6.3.1.1.4.1.1": v1.Mib{
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

// Test SNMP trapVars.
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
