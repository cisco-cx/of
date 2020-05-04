package v2_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	of "github.com/cisco-cx/of/pkg/v2"
	mib_registry "github.com/cisco-cx/of/wrap/mib/v2"
)

// Sample configs.
var YamlContent = `epc:
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
    - type: set
      key: alert_name
      value: starCard
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
    device_identifiers:
    - user-sha-aes128
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
        - .1.3.6.1.4.1.8164.2.13  # dummyTestValue
      annotation_mods: []  # this is allowed
device_not_found:
  defaults:
    source_type: cluster   # (host|cluster)...
      # if cluster, you must define defaults.clusters.
    clusters:
      nso1.example.org:  # cluster_name
        source_addresses:  # TODO: enhance this when necessary.
        - 192.168.1.28
        - dead:beef::1
    device_identifiers:
    - user-not-found
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
        - .1.3.6.1.4.1.8164.2.13  # dummyTestValue
      annotation_mods: []  # this is allowed`

var testEvents = `[{
    "apiVersion": "v1alpha1",
    "kind": "SNMPTrap",
    "receipts": {
      "snmptrapd": {
        "source": {
          "hostname": "localhost",
          "transportLayerProtocol": "UDP",
          "address": "::1",
          "port": "48381",
          "internetLayerProtocol": "IPv6"
        },
        "timestamp": "2019-04-26T03:46:57Z",
        "pduSecurity": "TRAP2, SNMP v3, user user-sha-aes128, context",
        "vars": [
          {
            "oid": ".1.3.6.1.2.1.1.3.0",
            "type": "Timeticks",
            "value": "(123) 0:00:01.23"
          },
          {
            "oid": ".1.3.6.1.6.3.1.1.4.1.0",
            "type": "OID",
            "value": ".1.3.6.1.4.1.8164.2.44"
          },
          {
            "oid": ".1.3.6.1.4.1.65000.1.1.1.1.1",
            "type": "STRING",
            "value": "foo"
          },
          {
            "oid": ".1.3.6.1.4.1.65000.1.1.1.1.1",
            "type": "STRING",
            "value": "bar"
          }
        ]
      },
      "filebeat": {
        "agent": {
          "id": "c9da5463-1d21-405e-afe8-c3aa1b7d3bba",
          "hostname": "vagrant",
          "ephemeral_id": "543bb87e-17ca-4aa3-a9fb-cc67aa28c5c4",
          "type": "filebeat",
          "version": "7.0.0"
        },
        "ecs": {
          "version": "1.0.0"
        },
        "log": {
          "offset": 0,
          "file": {
            "path": ""
          }
        },
        "@timestamp": "2019-04-26T03:46:57.941Z",
        "@version": "1",
        "message": "SNMPTRAP timestamp=[2019-04-26T03:46:57Z] hostname=[localhost] address=[UDP/IPv6: [::1]:48381] pdu_security=[TRAP2, SNMP v3, user user-sha-aes128, context ] vars[.1.3.6.1.2.1.1.3.0 = Timeticks: (123) 0:00:01.23\t.1.3.6.1.6.3.1.1.4.1.0 = OID: .1.3.6.1.6.3.1.1.5.1\t.1.3.6.1.4.1.65000.1.1.1.1.1 = STRING: \"foo\"\t.1.3.6.1.4.1.65000.1.1.1.1.1 = STRING: \"bar\"]",
        "input": {
          "type": "stdin"
        },
        "host": {
          "name": "vagrant"
        }
      },
      "logstash": {
        "tags": [
          "beats_input_codec_plain_applied"
        ]
      }
    }
  },{
    "apiVersion": "v1alpha1",
    "kind": "SNMPTrap",
    "receipts": {
      "snmptrapd": {
        "source": {
          "hostname": "localhost",
          "transportLayerProtocol": "UDP",
          "address": "::1",
          "port": "48381",
          "internetLayerProtocol": "IPv6"
        },
        "timestamp": "2019-04-26T03:46:57Z",
        "pduSecurity": "TRAP2, SNMP v3, user user-sha-aes128, context",
        "vars": [
          {
            "oid": ".1.3.6.1.2.1.1.3.0",
            "type": "Timeticks",
            "value": "(123) 0:00:01.23"
          },
          {
            "oid": ".1.3.6.1.6.3.1.1.4.1.0",
            "type": "OID",
            "value": ".1.3.6.1.4.1.8164.2.43"
          },
          {
            "oid": ".1.3.6.1.4.1.65000.1.1.1.1.1",
            "type": "STRING",
            "value": "foo"
          },
          {
            "oid": ".1.3.6.1.4.1.65000.1.1.1.1.1",
            "type": "STRING",
            "value": "bar"
          }
        ]
      },
      "filebeat": {
        "agent": {
          "id": "c9da5463-1d21-405e-afe8-c3aa1b7d3bba",
          "hostname": "vagrant",
          "ephemeral_id": "543bb87e-17ca-4aa3-a9fb-cc67aa28c5c4",
          "type": "filebeat",
          "version": "7.0.0"
        },
        "ecs": {
          "version": "1.0.0"
        },
        "log": {
          "offset": 0,
          "file": {
            "path": ""
          }
        },
        "@timestamp": "2019-04-26T03:46:57.941Z",
        "@version": "1",
        "message": "SNMPTRAP timestamp=[2019-04-26T03:46:57Z] hostname=[localhost] address=[UDP/IPv6: [::1]:48381] pdu_security=[TRAP2, SNMP v3, user user-sha-aes128, context ] vars[.1.3.6.1.2.1.1.3.0 = Timeticks: (123) 0:00:01.23\t.1.3.6.1.6.3.1.1.4.1.0 = OID: .1.3.6.1.6.3.1.1.5.1\t.1.3.6.1.4.1.65000.1.1.1.1.1 = STRING: \"foo\"\t.1.3.6.1.4.1.65000.1.1.1.1.1 = STRING: \"bar\"]",
        "input": {
          "type": "stdin"
        },
        "host": {
          "name": "vagrant"
        }
      },
      "logstash": {
        "tags": [
          "beats_input_codec_plain_applied"
        ]
      }
    }
  },{
    "apiVersion": "v1alpha1",
    "kind": "SNMPTrap",
    "receipts": {
      "snmptrapd": {
        "source": {
          "hostname": "localhost",
          "transportLayerProtocol": "UDP",
          "address": "::1",
          "port": "48381",
          "internetLayerProtocol": "IPv6"
        },
        "timestamp": "2019-04-26T03:46:57Z",
        "pduSecurity": "TRAP2, SNMP v3, user user-sha-aes128, context",
        "vars": [

		{
			"oid":   ".1.3.6.1.6.1.1.1.4.1",
			"value": ".1.3.6.1.4.1.8164.1.2.1.1.1"
		},
		{
			"oid":   ".1.3.6.1.4.1.8164.1.2.1.1.1",
			"value": "14"
		},
		{
			"oid":   ".1.3.6.1.4.1.24961.2.103.1.1.5.1.2",
			"value": "package-load-failure"
		},
		{
			"oid":   ".1.3.6.1.2.1.1.3.0",
			"type":  "Timeticks",
			"value": "(123) 0:00:01.23"
		},
		{
			"oid":   ".1.3.6.1.6.3.1.1.4.1",
			"type":  "OID",
			"value": ".1.3.6.1.4.1.8164.2.13"
		},
		{
			"oid":   ".1.3.6.1.6.3.1.1.4.1.0",
			"type":  "OID",
			"value": ".1.3.6.1.4.1.8164.2.44"
		},
		{
			"oid":   ".1.3.6.1.4.1.8164.2.44",
			"type":  "STRING",
			"value": "foo"
		},
		{
			"oid":   ".1.3.6.1.6.3.1.1.4.1.1",
			"type":  "OID",
			"value": ".1.3.6.1.4.1.8164.2.45"
		},
		{
			"oid":   ".1.3.6.1.4.1.8164.2.45",
			"type":  "OID",
			"value": ".1.3.6.1.4.1.65000.1.1.1.1.1"
		},
		{
			"oid":   ".1.3.6.1.4.1.65000.1.1.1.1.1",
			"type":  "STRING",
			"value": "bar"
		}
        ]
      },
      "filebeat": {
        "agent": {
          "id": "c9da5463-1d21-405e-afe8-c3aa1b7d3bba",
          "hostname": "vagrant",
          "ephemeral_id": "543bb87e-17ca-4aa3-a9fb-cc67aa28c5c4",
          "type": "filebeat",
          "version": "7.0.0"
        },
        "ecs": {
          "version": "1.0.0"
        },
        "log": {
          "offset": 0,
          "file": {
            "path": ""
          }
        },
        "@timestamp": "2019-04-26T03:46:57.941Z",
        "@version": "1",
        "message": "SNMPTRAP timestamp=[2019-04-26T03:46:57Z] hostname=[localhost] address=[UDP/IPv6: [::1]:48381] pdu_security=[TRAP2, SNMP v3, user user-sha-aes128, context ] vars[.1.3.6.1.2.1.1.3.0 = Timeticks: (123) 0:00:01.23\t.1.3.6.1.6.3.1.1.4.1.0 = OID: .1.3.6.1.6.3.1.1.5.1\t.1.3.6.1.4.1.65000.1.1.1.1.1 = STRING: \"foo\"\t.1.3.6.1.4.1.65000.1.1.1.1.1 = STRING: \"bar\"]",
        "input": {
          "type": "stdin"
        },
        "host": {
          "name": "vagrant"
        }
      },
      "logstash": {
        "tags": [
          "beats_input_codec_plain_applied"
        ]
      }
    }
  }]`

// Test mibs data.
func mibRegistry(t *testing.T) of.MIBRegistry {
	mibs := map[string]of.MIB{
		".1.3.6.1.2.1.1.3.0": of.MIB{
			Name: "oid1",
		},
		".1.3.6.1.6.3.1.1.4.1.0": of.MIB{
			Name: "oid2",
		},
		".1.3.6.1.4.1.8164.2.44": of.MIB{
			Name: "oid3",
		},
		".1.3.6.1.4.1.8164.2.45": of.MIB{
			Name: "oid4",
		},
		".1.3.6.1.4.1.65000.1.1.1.1.1": of.MIB{
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
			Oid:   ".1.3.6.1.6.1.1.1.4.1",
			Value: ".1.3.6.1.4.1.8164.1.2.1.1.1",
		},
		of.TrapVar{
			Oid:   ".1.3.6.1.4.1.8164.1.2.1.1.1",
			Value: "14",
		},
		of.TrapVar{
			Oid:   ".1.3.6.1.4.1.24961.2.103.1.1.5.1.2",
			Value: "package-load-failure",
		},
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

func TrapSource() *of.TrapSource {
	return &of.TrapSource{
		Address:  "192.168.1.28",
		Hostname: "localhost",
	}
}

func TrapReceipts() *of.Receipts {
	return &of.Receipts{
		Snmptrapd: of.Snmptrapd{
			Timestamp:   "2019-04-26T03:46:57Z",
			Source:      *TrapSource(),
			Vars:        *trapVars(),
			PduSecurity: "TRAP2, SNMP v3, user user-sha-aes128, context",
		},
		Filebeat: of.Filebeat{
			Message:   "SNMPTRAP timestamp=[2019-04-26T03:46:57Z] hostname=[localhost] address=[UDP/IPv6: [::1]:48381] pdu_security=[TRAP2, SNMP v3, user user-sha-aes128, context ] vars[.1.3.6.1.2.1.1.3.0 = Timeticks: (123) 0:00:01.23\t.1.3.6.1.6.3.1.1.4.1.0 = OID: .1.3.6.1.6.3.1.1.5.1\t.1.3.6.1.4.1.65000.1.1.1.1.1 = STRING: \"foo\"\t.1.3.6.1.4.1.65000.1.1.1.1.1 = STRING: \"bar\"]",
			Timestamp: "2019-04-26T03:46:57.941Z",
		},
	}
}

func TrapEvents() []*of.PostableEvent {
	events := []*of.PostableEvent{
		&of.PostableEvent{
			Document: of.Document{
				Receipts: *TrapReceipts(),
			},
		},
	}
	return events
}
