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

package v1

import (
	"fmt"
	"strings"

	of "github.com/cisco-cx/of/pkg/v1"
)

type MibRegistry struct {
	regs  map[string]of.Mib
	index map[string][]string
}

func New() *MibRegistry {
	regs := make(map[string]of.Mib)
	index := make(map[string][]string)
	return &MibRegistry{
		regs:  regs,
		index: index,
	}
}

// Return Mib for given OID.
func (mib *MibRegistry) Mib(oid string) of.Mib {
	return mib.regs[oid]
}

// Return the last node to its name. Ex: 1.3.6.1.2.1.11.19 -> snmpInTraps.
func (mib *MibRegistry) ShortString(oid string) string {
	return mib.regs[oid].Name
}

// Return display string for given OID.
func (mib *MibRegistry) String(oid string) string {
	if value, hasValue := mib.index[oid]; hasValue == true {
		return strings.Join(value, ".")
	}
	mib.index[oid] = mib.getStrOid(oid)
	return mib.String(oid)
}

func (mib *MibRegistry) getStrOid(oid string) []string {
	mibReg := mib.Mib(oid)
	idx := strings.LastIndex(oid, ".")
	if idx == -1 {
		// if first oid element is not loaded then empty string to display
		if len(mibReg.Name) <= 0 {
			return []string{}
		}
		return []string{mibReg.Name}
	}

	strOid := mib.getStrOid(oid[:idx])
	if len(strOid) <= 0 {
		return []string{}
	}

	// if not loaded but previous elements then we show the best effort
	// e.g. 1st-Name.2nd-Name.5.6. See tests for more examples
	name := mibReg.Name
	if len(name) <= 0 {
		name = oid[idx+1:]
	}
	return append(strOid, name)
}

// Load given map[oid]Mib into registry.
func (mib *MibRegistry) Load(src map[string]of.Mib) error {
	for k, v := range src {
		if len(v.Name) <= 0 {
			return of.Error(fmt.Sprintf("Name can't be empty: '%+v'", v))
		}
		mib.regs[k] = v
	}
	return nil
}
