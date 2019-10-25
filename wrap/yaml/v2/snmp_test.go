// Copyright 2019 Cisco Systems, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package v1_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	of "github.com/cisco-cx/of/pkg/v1"
	snmp_config "github.com/cisco-cx/of/pkg/v2/snmp"
	yaml "github.com/cisco-cx/of/wrap/yaml/v2"
)

// Sample EPC config with comments removed for comparison.
var yamlContent = `epc:
  defaults:
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
        - .1.3.6.1.4.1.8164.2.55`

// Expected Golang structure for above config.
var expectedCfg = yaml.Configs{
	"epc": snmp_config.Config{
		Defaults: snmp_config.Default{
			Enabled:            nil,
			SourceType:         snmp_config.HostType,
			GeneratorUrlPrefix: "http://www.oid-info.com/get/",
			LabelMods: []snmp_config.Mod{
				snmp_config.Mod{
					Type:  snmp_config.Set,
					Key:   "vendor",
					Value: "cisco",
				},
				snmp_config.Mod{
					Type:  snmp_config.Set,
					Key:   "subsystem",
					Value: "epc",
				},
				snmp_config.Mod{
					Type:    snmp_config.Copy,
					Oid:     ".1.3.6.1.4.1.8164.1.2.1.1.1",
					As:      snmp_config.Value,
					ToKey:   "star_slot_num",
					OnError: snmp_config.Drop,
				},
			},
			AnnotationMods: []snmp_config.Mod{

				snmp_config.Mod{
					Type:  snmp_config.Copy,
					Oid:   ".1.3.6.1.6.1.1.1.4.1",
					As:    snmp_config.Value,
					ToKey: "event_oid",
				},
				snmp_config.Mod{
					Type:  snmp_config.Copy,
					Oid:   ".1.3.6.1.6.1.1.1.4.1",
					As:    snmp_config.OidValueStrShort,
					ToKey: "event_name",
				},
			},
		},
		Alerts: []snmp_config.Alert{
			snmp_config.Alert{
				Name:    "starCard",
				Enabled: func() *bool { b := true; return &b }(),
				LabelMods: []snmp_config.Mod{
					snmp_config.Mod{
						Type:  snmp_config.Set,
						Key:   "alert_severity",
						Value: "error",
					},
				},
				Firing: map[string][]snmp_config.Select{
					"select": []snmp_config.Select{
						snmp_config.Select{
							Type: snmp_config.Equals,
							Oid:  ".1.3.6.1.6.3.1.1.4.1",
							As:   snmp_config.Value,
							Values: []string{
								".1.3.6.1.4.1.8164.2.13",
								".1.3.6.1.4.1.8164.2.4",
								".1.3.6.1.4.1.8164.2.7",
							},
						},
					},
				},
				Clearing: map[string][]snmp_config.Select{
					"select": []snmp_config.Select{
						snmp_config.Select{
							Type: snmp_config.Equals,
							Oid:  ".1.3.6.1.6.3.1.1.4.1",
							As:   snmp_config.Value,
							Values: []string{
								".1.3.6.1.4.1.8164.2.5",
								".1.3.6.1.4.1.8164.2.55",
							},
						},
					},
				},
			},
			snmp_config.Alert{
				Name: "starCardBootFailed",
				LabelMods: []snmp_config.Mod{
					snmp_config.Mod{
						Type:  snmp_config.Set,
						Key:   "alert_severity",
						Value: "error",
					},
				},
				Firing: map[string][]snmp_config.Select{
					"select": []snmp_config.Select{
						snmp_config.Select{
							Type: snmp_config.Equals,
							Oid:  ".1.3.6.1.6.3.1.1.4.1",
							As:   snmp_config.Value,
							Values: []string{
								".1.3.6.1.4.1.8164.2.9",
							},
						},
					},
				},
				Clearing: map[string][]snmp_config.Select{
					"select": []snmp_config.Select{
						snmp_config.Select{
							Type: snmp_config.Equals,
							Oid:  ".1.3.6.1.6.3.1.1.4.1",
							As:   snmp_config.Value,
							Values: []string{
								".1.3.6.1.4.1.8164.2.5",
								".1.3.6.1.4.1.8164.2.55",
							},
						},
					},
				},
			},
			snmp_config.Alert{
				Name: "starCardActive",
				LabelMods: []snmp_config.Mod{
					snmp_config.Mod{
						Type:  snmp_config.Set,
						Key:   "alert_severity",
						Value: "informational",
					},
				},
				Firing: map[string][]snmp_config.Select{
					"select": []snmp_config.Select{
						snmp_config.Select{
							Type: snmp_config.Equals,
							Oid:  ".1.3.6.1.6.3.1.1.4.1",
							As:   snmp_config.Value,
							Values: []string{
								".1.3.6.1.4.1.8164.2.55",
							},
						},
					},
				},
			},
		},
	},
}

// Enforce interface implementation.
func TestAlertsInterface(t *testing.T) {
	var _ of.Decoder = &yaml.Configs{}
	var _ of.Encoder = &yaml.Configs{}
}

// Ensure yaml decodes Alerts
func TestAlertsDecoder(t *testing.T) {

	r := strings.NewReader(yamlContent)

	cfg := yaml.Configs{}
	err := cfg.Decode(r)
	require.NoError(t, err)
	require.EqualValues(t, expectedCfg, cfg)
}

// Ensure yaml encodes Alerts
func TestAlertsEncoder(t *testing.T) {

	cfg := yaml.Configs{}
	buf := bytes.NewBuffer(nil)
	cfg["epc"] = expectedCfg["epc"]
	err := cfg.Encode(buf)
	require.NoError(t, err)
	//fmt.Printf("%+v\n", string(buf.Bytes()))
	require.EqualValues(t, strings.Trim(yamlContent, "\n"), strings.Trim(string(buf.Bytes()), "\n"))
}
