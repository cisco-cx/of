package cmd_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	of "github.com/cisco-cx/of/pkg/v2"
	mib_registry "github.com/cisco-cx/of/wrap/mib/v2"
)

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
			Timestamp: "2019-04-26T03:46:57Z",
			Source:    *TrapSource(),
			Vars:      *trapVars(),
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
