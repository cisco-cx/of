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
	"strconv"

	of "github.com/cisco-cx/of/pkg/v1"
)

// ACIFaultSeverityID represents a ACI fault's severity level ID as mapped
// inside the Observability Framework.
//
// ACIFaultSeverityID implements the of.ACIFaultSeverityIDParser interface.
type ACIFaultSeverityID struct {
	ofID of.ACIFaultSeverityID
}

// NewACIFaultSeverityID returns a new instance of ACIFaultSeverityID.
func NewACIFaultSeverityID(id int) (ACIFaultSeverityID, error) {
	return ACIFaultSeverityID{
		ofID: of.ACIFaultSeverityID(id),
	}, nil
}

// Raw returns the ACIFaultSeverityID's equivalent ACIFaultSeverityRaw.
func (s ACIFaultSeverityID) Raw() of.ACIFaultSeverityRaw {
	m := map[int]string{
		0: "cleared",
		1: "info",
		2: "warning",
		3: "minor",
		4: "major",
		5: "critical",
	}
	r := m[int(s.ofID)]
	return of.ACIFaultSeverityRaw(r)
}

// String implements the fmt.Stringer interface.
func (s ACIFaultSeverityID) String() string {
	return strconv.Itoa(int(s.ofID))
}
