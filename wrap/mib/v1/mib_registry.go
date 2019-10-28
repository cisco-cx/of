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

// MIBRegistry keeps the data for building the (OID, strings) map
type MIBRegistry struct {
	regs  map[string]*of.MIB
	index map[string][]string
}

// Return a new MIBRegistry pointer
func New() *MIBRegistry {
	regs := make(map[string]*of.MIB)
	index := make(map[string][]string)
	return &MIBRegistry{
		regs:  regs,
		index: index,
	}
}

// Return MIB for given OID.
func (MIB *MIBRegistry) MIB(oid string) *of.MIB {
	return MIB.regs[oid]
}

// Return the last node to its name. Ex: 1.3.6.1.2.1.11.19 -> snmpInTraps.
func (MIB *MIBRegistry) ShortString(oid string) string {
	if r := MIB.MIB(oid); r != nil {
		return r.Name
	}
	return ""
}

// Return display string for given OID.
func (MIB *MIBRegistry) String(oid string) string {
	if value, hasValue := MIB.index[oid]; hasValue == true {
		return strings.Join(value, ".")
	}
	MIB.index[oid] = MIB.getStrOID(oid)
	return MIB.String(oid)
}

func (MIB *MIBRegistry) getStrOID(oid string) []string {
	mibReg := MIB.MIB(oid)
	if mibReg != nil {
		idx := strings.LastIndex(oid, ".")
		if idx == -1 {
			return []string{mibReg.Name}
		}

		strOid := MIB.getStrOID(oid[:idx])
		return append(strOid, mibReg.Name)
	}

	idx := strings.LastIndex(oid, ".")
	if idx == -1 {
		return []string{oid}
	}

	strOid := MIB.getStrOID(oid[:idx])
	return append(strOid, oid[idx+1:])
}

// Load given map[oid]MIB into registry.
func (MIB *MIBRegistry) Load(src map[string]of.MIB) error {
	for k, v := range src {
		if len(v.Name) <= 0 {
			return of.Error(fmt.Sprintf("Name can't be empty: '%+v'", v))
		}
		v_copy_ptr := new(of.MIB)
		*v_copy_ptr = v
		MIB.regs[k] = v_copy_ptr
	}
	return nil
}
