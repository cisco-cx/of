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

// MibRegistry keeps the data for building the (OID, strings) map
type MibRegistry struct {
	regs  map[string]*of.Mib
	index map[string][]string
}

// Return a new MibRegistry pointer
func New() *MibRegistry {
	regs := make(map[string]*of.Mib)
	index := make(map[string][]string)
	return &MibRegistry{
		regs:  regs,
		index: index,
	}
}

// Return Mib for given OID.
func (mib *MibRegistry) Mib(oid string) *of.Mib {
	return mib.regs[oid]
}

// Return the last node to its name. Ex: 1.3.6.1.2.1.11.19 -> snmpInTraps.
func (mib *MibRegistry) ShortString(oid string) string {
	if r := mib.Mib(oid); r != nil {
		return r.Name
	}
	return ""
}

// Return display string for given OID.
func (mib *MibRegistry) String(oid string) string {
	if value, hasValue := mib.index[oid]; hasValue == true {
		return strings.Join(value, ".")
	}
	mib.index[oid] = mib.getStrOID(oid)
	return mib.String(oid)
}

func (mib *MibRegistry) getStrOID(oid string) []string {
	mibReg := mib.Mib(oid)
	if mibReg != nil {
		idx := strings.LastIndex(oid, ".")
		if idx == -1 {
			return []string{mibReg.Name}
		}

		strOid := mib.getStrOID(oid[:idx])
		return append(strOid, mibReg.Name)
	}

	idx := strings.LastIndex(oid, ".")
	if idx == -1 {
		return []string{oid}
	}

	strOid := mib.getStrOID(oid[:idx])
	return append(strOid, oid[idx+1:])
}

// Load given map[oid]Mib into registry.
func (mib *MibRegistry) Load(src map[string]of.Mib) error {
	for k, v := range src {
		if len(v.Name) <= 0 {
			return of.Error(fmt.Sprintf("Name can't be empty: '%+v'", v))
		}
		v_copy_ptr := new(of.Mib)
		*v_copy_ptr = v
		mib.regs[k] = v_copy_ptr
	}
	return nil
}
