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

package v2

// This package represents snmptrap data.

type Agent struct {
	Hostname    string `json:"hostname,omitempty"`
	ID          string `json:"id,omitempty"`
	EphemeralID string `json:"ephemeral_id,omitempty"`
	Type        string `json:"type,omitempty"`
	Version     string `json:"version,omitempty"`
}

type AppStatus struct {
	ApiVersion  string         `json:"apiVersion"`
	Description string         `json:"description"`
	Links       AppStatusLinks `json:"links"`
	Status      string         `json:"status"`
}

type AppStatusLinks struct {
	About string `json:"about"`
}

type Document struct {
	ApiVersion string   `json:"apiVersion,omitempty"`
	Kind       string   `json:"kind,omitempty"`
	Receipts   Receipts `json:"receipts,omitempty"`
}

type Ecs struct {
	Version string `json:"version,omitempty"`
}

type Filebeat struct {
	Agent     Agent              `json:"agent,omitempty"`
	Ecs       Ecs                `json:"ecs,omitempty"`
	Input     PostableEventInput `json:"input,omitempty"`
	Host      Host               `json:"host,omitempty"`
	Log       Log                `json:"log,omitempty"`
	Message   string             `json:"message,omitempty"`
	Version   string             `json:"@version,omitempty"`
	Timestamp string             `json:"@timestamp,omitempty"`
}

type Host struct {
	Name string `json:"name,omitempty"`
}

type Log struct {
	Offset int32   `json:"offset,omitempty"`
	File   LogFile `json:"file,omitempty"`
}

type LogFile struct {
	Path string `json:"path,omitempty"`
}

type Logstash struct {
	Tags []string `json:"tags,omitempty"`
}

type PostableEvent struct {
	Agent      Agent              `json:"agent,omitempty"`
	ApiVersion string             `json:"apiVersion,omitempty"`
	Document   Document           `json:"document,omitempty"`
	Ecs        Ecs                `json:"ecs,omitempty"`
	Input      PostableEventInput `json:"input,omitempty"`
	Host       Host               `json:"host,omitempty"`
	Log        Log                `json:"log,omitempty"`
	Message    string             `json:"message,omitempty"`
	Version    string             `json:"@version,omitempty"`
	Tags       []string           `json:"tags,omitempty"`
	Timestamp  string             `json:"@timestamp,omitempty"`
}

type PostableEventInput struct {
	Type string `json:"type,omitempty"`
}

type Receipts struct {
	Filebeat  Filebeat  `json:"filebeat,omitempty"`
	Logstash  Logstash  `json:"logstash,omitempty"`
	Snmptrapd Snmptrapd `json:"snmptrapd,omitempty"`
}

type Snmptrapd struct {
	Timestamp   string     `json:"timestamp,omitempty"`
	Source      TrapSource `json:"source,omitempty"`
	Vars        []TrapVar  `json:"vars,omitempty"`
	PduSecurity string     `json:"pduSecurity,omitempty"`
}

type TrapSource struct {
	Address                string `json:"address"`
	Hostname               string `json:"hostname"`
	InternetLayerProtocol  string `json:"internetLayerProtocol"`
	Port                   string `json:"port"`
	TransportLayerProtocol string `json:"transportLayerProtocol"`
}

type TrapVar struct {
	Value string `json:"value"`
	Type  string `json:"type"`
	Oid   string `json:"oid"`
}
