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

import "fmt"

// ACIFaultSeverityRaw represents a ACI fault's raw severity level as returned by
// the ACI API.
type ACIFaultSeverityRaw string

// ACIFaultSeverityRawParser represents the ability, given a ACIFaultSeverityRaw,
// to parse it into different formats (e.g. ACIFaultSeverityID).
type ACIFaultSeverityRawParser interface {
	fmt.Stringer
	ID() ACIFaultSeverityID
}
