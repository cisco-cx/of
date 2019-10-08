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
//
// This work incorporates works covered by the following notices:
//

package v1_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	of "github.com/cisco-cx/of/lib/v1"
	acigo "github.com/cisco-cx/of/wrap/acigo/v1"
)

// t1 is a private struct that contains all valid {ID, Raw} values.
var t1 = []struct {
	id  int
	raw string
}{
	{0, "cleared"},
	{1, "info"},
	{2, "warning"},
	{3, "minor"},
	{4, "major"},
	{5, "critical"},
}

// Confirm that acigo.ACIFaultSeverityID implements the
// of.ACIFaultSeverityIDParser interface.
func TestACIFaultSeverityID_InterfaceACIFaultSeverityIDParser(t *testing.T) {
	var _ of.ACIFaultSeverityIDParser = acigo.ACIFaultSeverityID{}
	assert.Nil(t, nil) // If we get this far, the test passed.
}

// NewACIFaultSeverityID() Confirm simple functionality.
func TestACIFaultSeverityID_Simple(t *testing.T) {
	_, err := acigo.NewACIFaultSeverityID(4)
	assert.Nil(t, err)
}

// ACIFaultSeverityID.Raw() Test Positive-path functionality.
func TestACIFaultSeverityID_RawPositive(t *testing.T) {
	for _, tr := range t1 {
		t.Run(string(tr.id), func(t *testing.T) {
			// Prepare to assert multiple times.
			assert := assert.New(t)
			// Pass in ID from table to make new instance.
			id, err := acigo.NewACIFaultSeverityID(tr.id)
			assert.Nil(err)
			assert.Equal(tr.raw, string(id.Raw()))
		})
	}
}

// Simple test of ACIFaultSeverityID's implementation of the fmt.Stringer
// interface.
func TestACIFaultSeverityID_StringerSimple(t *testing.T) {
	// Prepare to assert multiple times.
	assert := assert.New(t)
	id, err := acigo.NewACIFaultSeverityID(5)
	assert.Nil(err)
	assert.Equal("5", fmt.Sprintf("%s", id))
}
