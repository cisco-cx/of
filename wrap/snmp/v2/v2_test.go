package v2_test

import (
	"fmt"
	"testing"

	of "github.com/cisco-cx/of/pkg/v2"
)

// Sample configs.
var YamlConfigs = `config1:
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
      value: config1
    annotation_mods:
    - type: set
      key: vendor
      value: cisco
    - type: set
      key: subsystem
      value: config1
  alerts:
  - name: starTaskFailure
    label_mods:
    - type: set
      key: alertname
      value: starTaskFailure
    - type: set
      key: alert_severity
      value: major
    firing:
      select:
      - type: equals
        oid: .1.3.6.1.6.3.1.1.4.1.0
        as: value
        values:
        - .1.3.6.1.4.1.8164.2.150
        annotation_mods:
        - type: copy
          oid: .1.3.6.1.6.3.1.1.4.1.0
          as: value
          to_key: event_name
          map:
            .1.3.6.1.4.1.8164.2.150: starTaskFailed
        - type: set
          key: compatible_clear_events
          value: '{".1.3.6.1.4.1.8164.2.151":{"event_name":"starTaskRestart"}}'
    clearing:
      select:
      - type: equals
        oid: .1.3.6.1.6.3.1.1.4.1.0
        as: value
        values:
        - .1.3.6.1.4.1.8164.2.151
        annotation_mods:
        - type: copy
          oid: .1.3.6.1.6.3.1.1.4.1.0
          as: value
          to_key: event_name
          map:
            .1.3.6.1.4.1.8164.2.151: starTaskRestart
  - name: starTaskRestart
    label_mods:
    - type: set
      key: alertname
      value: starTaskRestart
    - type: set
      key: alert_severity
      value: major
    firing:
      select:
      - type: equals
        oid: .1.3.6.1.6.3.1.1.4.1.0
        as: value
        values:
        - .1.3.6.1.4.1.8164.2.150
        - .1.3.6.1.4.1.8164.2.151
        annotation_mods:
        - type: copy
          oid: .1.3.6.1.6.3.1.1.4.1.0
          as: value
          to_key: event_name
          map:
            .1.3.6.1.4.1.8164.2.150: starTaskFailed
            .1.3.6.1.4.1.8164.2.151: starTaskRestart`

var StarEvents = `[{
  "apiVersion": "v1alpha1",
  "kind": "SNMPTrap",
  "receipts": {
    "filebeat": {
      "@timestamp": "2020-05-01T22:22:53.324Z",
      "@version": "1",
      "agent": {
        "ephemeral_id": "727040a1-9f47-4a94-8b32-47636376ce45",
        "hostname": "test-host-01",
        "id": "1c0e1fe8-3c78-4fe1-896d-ed1d2342f12c",
        "type": "filebeat",
        "version": "7.2.0"
      },
      "ecs": {
        "version": "1.0.0"
      },
      "host": {
        "name": "test-host-01"
      },
      "input": {
        "type": "stdin"
      },
      "log": {
        "file": {
          "path": ""
        },
        "offset": 0
      },
      "message": "SNMPTRAP timestamp=[2020-05-01T22:22:53Z] address=[UDP/IPv6: [dead::beef]:46819] pdu_security=[TRAP2, SNMP v3, user snmp-user, context ] vars[.1.3.6.1.2.1.1.3.0 = Timeticks: (290240897) 33 days, 14:13:28.97\t.1.3.6.1.6.3.1.1.4.1.0 = OID: .1.3.6.1.4.1.8164.2.150\t.1.3.6.1.4.1.8164.1.20.1.1.3 = STRING: ;sessmgr;\t.1.3.6.1.4.1.8164.1.20.1.1.2 = Gauge32: 34\t.1.3.6.1.4.1.8164.1.20.1.1.4 = Gauge32: 3\t.1.3.6.1.4.1.8164.1.20.1.1.5 = Gauge32: 0\t.1.3.6.1.6.3.1.1.4.3.0 = OID: .1.3.6.1.4.1.8164.2]"
    },
    "logstash": {
      "tags": [
        "beats_input_codec_plain_applied"
      ]
    },
    "snmptrapd": {
      "pduSecurity": "TRAP2, SNMP v3, user snmp-user, context",
      "source": {
        "address": "dead::beef",
        "hostname": "test-device-01",
        "internetLayerProtocol": "IPv6",
        "port": "46819",
        "transportLayerProtocol": "UDP"
      },
      "timestamp": "2020-05-01T22:22:53Z",
      "vars": [
        {
          "oid": ".1.3.6.1.2.1.1.3.0",
          "type": "Timeticks",
          "value": "(290240897) 33 days, 14:13:28.97"
        },
        {
          "oid": ".1.3.6.1.6.3.1.1.4.1.0",
          "type": "OID",
          "value": ".1.3.6.1.4.1.8164.2.150"
        },
        {
          "oid": ".1.3.6.1.4.1.8164.1.20.1.1.3",
          "type": "STRING",
          "value": ";sessmgr;"
        },
        {
          "oid": ".1.3.6.1.4.1.8164.1.20.1.1.2",
          "type": "Gauge32",
          "value": "34"
        },
        {
          "oid": ".1.3.6.1.4.1.8164.1.20.1.1.4",
          "type": "Gauge32",
          "value": "3"
        },
        {
          "oid": ".1.3.6.1.4.1.8164.1.20.1.1.5",
          "type": "Gauge32",
          "value": "0"
        },
        {
          "oid": ".1.3.6.1.6.3.1.1.4.3.0",
          "type": "OID",
          "value": ".1.3.6.1.4.1.8164.2"
        }
      ]
    }
  }
},
{
  "apiVersion": "v1alpha1",
  "kind": "SNMPTrap",
  "receipts": {
    "filebeat": {
      "@timestamp": "2020-05-01T22:22:53.324Z",
      "@version": "1",
      "agent": {
        "ephemeral_id": "e45708c5-6d8e-4f33-9a68-c049124d6685",
        "hostname": "test-host-01",
        "id": "41778bee-763a-4b3b-8cd4-cb81f7788e49",
        "type": "filebeat",
        "version": "7.2.0"
      },
      "ecs": {
        "version": "1.0.0"
      },
      "host": {
        "name": "test-host-01"
      },
      "input": {
        "type": "stdin"
      },
      "log": {
        "file": {
          "path": ""
        },
        "offset": 0
      },
      "message": "SNMPTRAP timestamp=[2020-05-01T22:23:53Z] address=[UDP/IPv6: [dead::beef]:53183] pdu_security=[TRAP2, SNMP v3, user snmp-user, context ] vars[.1.3.6.1.2.1.1.3.0 = Timeticks: (290240897) 33 days, 14:13:28.97\t.1.3.6.1.6.3.1.1.4.1.0 = OID: .1.3.6.1.4.1.8164.2.151\t.1.3.6.1.4.1.8164.1.20.1.1.3 = STRING: ;sessmgr;\t.1.3.6.1.4.1.8164.1.20.1.1.2 = Gauge32: 34\t.1.3.6.1.4.1.8164.1.20.1.1.4 = Gauge32: 3\t.1.3.6.1.4.1.8164.1.20.1.1.5 = Gauge32: 0\t.1.3.6.1.6.3.1.1.4.3.0 = OID: .1.3.6.1.4.1.8164.2]"
    },
    "logstash": {
      "tags": [
        "beats_input_codec_plain_applied"
      ]
    },
    "snmptrapd": {
      "pduSecurity": "TRAP2, SNMP v3, user snmp-user, context",
      "source": {
        "address": "dead::beef",
        "hostname": "test-device-01",
        "internetLayerProtocol": "IPv6",
        "port": "53183",
        "transportLayerProtocol": "UDP"
      },
      "timestamp": "2020-05-01T22:23:53Z",
      "vars": [
        {
          "oid": ".1.3.6.1.2.1.1.3.0",
          "type": "Timeticks",
          "value": "(290240897) 33 days, 14:13:28.97"
        },
        {
          "oid": ".1.3.6.1.6.3.1.1.4.1.0",
          "type": "OID",
          "value": ".1.3.6.1.4.1.8164.2.151"
        },
        {
          "oid": ".1.3.6.1.4.1.8164.1.20.1.1.3",
          "type": "STRING",
          "value": ";sessmgr;"
        },
        {
          "oid": ".1.3.6.1.4.1.8164.1.20.1.1.2",
          "type": "Gauge32",
          "value": "34"
        },
        {
          "oid": ".1.3.6.1.4.1.8164.1.20.1.1.4",
          "type": "Gauge32",
          "value": "3"
        },
        {
          "oid": ".1.3.6.1.4.1.8164.1.20.1.1.5",
          "type": "Gauge32",
          "value": "0"
        },
        {
          "oid": ".1.3.6.1.6.3.1.1.4.3.0",
          "type": "OID",
          "value": ".1.3.6.1.4.1.8164.2"
        }
      ]
    }
  }
}
]`

type fakeMibRegistry struct {
}

func (f *fakeMibRegistry) MIB(string) *of.MIB {
	return &of.MIB{}
}

func (f *fakeMibRegistry) String(oid string) string {
	return oid
}

func (f *fakeMibRegistry) ShortString(oid string) string {
	return fmt.Sprintf("short_%s", oid)
}

func (f *fakeMibRegistry) Load(mibs map[string]of.MIB) error {
	return nil
}

func TestMibsInterface(t *testing.T) {
	var _ of.MIBRegistry = &fakeMibRegistry{}
}

func newFakeMibRegistry() of.MIBRegistry {
	return &fakeMibRegistry{}
}
