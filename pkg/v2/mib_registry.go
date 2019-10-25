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

// Represents Mib description.
type Mib struct {
	Name        string
	Description string
	Units       string
}

type MibRegistry interface {
	// Return Mib for given OID.
	Mib(string) Mib

	// Return Mib for given OID.
	// Translate each node in OID to its corresponding name, if MIB has its definition, else use the number.
	// Ex : 1.3.6.1.2.1.11.19 -> iso.org.dod.internet.mgmt.mib-2.snmp.snmpInTraps.
	//      1.3.6.1.2.1.11.19.54334 -> iso.org.dod.internet.mgmt.mib-2.snmp.snmpInTraps.54334.
	String(string) string

	// Translate the last node to its name. Ex: 1.3.6.1.2.1.11.19 -> snmpInTraps.
	ShortString(string) string

	// Load given map[oid]Mib into registry.
	Load(map[string]Mib) error
}
