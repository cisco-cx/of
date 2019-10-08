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
	of "github.com/cisco-cx/of/lib/v1"
)

// ACIFaultSeverityRaw represents a ACI fault's severity level string as mapped
// inside the Observability Framework.
//
// ACIFaultSeverityRaw implements the of.ACIFaultSeverityRawParser interface.
type ACIFaultSeverityRaw struct {
	ofRaw of.ACIFaultSeverityRaw
}

// NewACIFaultSeverityRaw returns a new instance of ACIFaultSeverityRaw.
func NewACIFaultSeverityRaw(s string) (ACIFaultSeverityRaw, error) {
	return ACIFaultSeverityRaw{
		ofRaw: of.ACIFaultSeverityRaw(s),
	}, nil
}

// ID returns the ACIFaultSeverityRaw's equivalent ACIFaultSeverityID.
func (s ACIFaultSeverityRaw) ID() of.ACIFaultSeverityID {
	m := map[string]int{
		"cleared":  0,
		"info":     1,
		"warning":  2,
		"minor":    3,
		"major":    4,
		"critical": 5,
	}
	return of.ACIFaultSeverityID(m[string(s.ofRaw)])
}

// String implements the fmt.Stringer interface.
func (s ACIFaultSeverityRaw) String() string {
	return fmt.Sprintf("%s", string(s.ofRaw))
}
