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

package aci

type Alerts struct {
	APIC          AlertsConfigAPIC                    `yaml:"apic,omitempty"`
	Defaults      AlertsConfigDefaults                `yaml:"defaults,omitempty"`
	DroppedFaults map[string]AlertsConfigDroppedFault `yaml:"dropped_faults,omitempty"`
	Alerts        map[string]AlertConfig              `yaml:"alerts,omitempty"`
}

type AlertsConfigAPIC struct {
	AlertSeverityThreshold string `yaml:"alert_severity_threshold,omitempty"`
	DropUnknownAlerts      bool   `yaml:"drop_unknown_alerts,omitempty"`
}

type AlertsConfigDefaults struct {
	AlertSeverity string `yaml:"alert_severity,omitempty"`
}

type AlertsConfigDroppedFault struct {
	FaultName string `yaml:"fault_name,omitempty"`
}

type AlertConfig struct {
	AlertSeverity string                      `yaml:"alert_severity,omitempty"`
	Faults        map[string]AlertConfigFault `yaml:"faults,omitempty"`
}

type AlertConfigFault struct {
	FaultName string `yaml:"fault_name,omitempty"`
}

type Secrets struct {
	APIC   SecretsConfigAPIC `yaml:"apic,omitempty"`
	Consul ConsulConfig      `yaml:"consul,omitempty"`
}

type SecretsConfigAPIC struct {
	Cluster SecretsConfigAPICCluster `yaml:"cluster,omitempty"`
}

type SecretsConfigAPICCluster struct {
	Name string `yaml:"name,omitempty"`
}

type ConsulConfig struct {
	Host     string            `yaml:"host"`
	Service  string            `yaml:"service"`
	NodeMeta map[string]string `yaml:"node_meta,omitempty"`
}
