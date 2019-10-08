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

// t2 is a private struct that contains all valid {ID, Raw} values.
var t2 = []struct {
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

// Confirm that acigo.ACIFaultSeverityRaw implements the
// of.ACIFaultSeverityRawParser interface.
func TestACIFaultSeverityRaw_InterfaceACIFaultSeverityRawParser(t *testing.T) {
	var _ of.ACIFaultSeverityRawParser = acigo.ACIFaultSeverityRaw{}
	assert.Nil(t, nil) // If we get this far, the test passed.
}

// NewACIFaultSeverityRaw() Confirm simple functionality.
func TestACIFaultSeverityRaw_Simple(t *testing.T) {
	_, err := acigo.NewACIFaultSeverityRaw("cleared")
	assert.Nil(t, err)
}

// ACIFaultSeverityRaw.ID() Test Positive-path functionality.
func TestACIFaultSeverityRaw_IDPositive(t *testing.T) {
	for _, tr := range t2 {
		t.Run(tr.raw, func(t *testing.T) {
			// Prepare to assert multiple times.
			assert := assert.New(t)
			// Pass in Raw from table to make new instance.
			raw, err := acigo.NewACIFaultSeverityRaw(tr.raw)
			assert.Nil(err)
			assert.Equal(tr.id, int(raw.ID()))
		})
	}
}

// Simple test of ACIFaultSeverityRaw's implementation of the fmt.Stringer
// interface.
func TestACIFaultSeverityRaw_StringerSimple(t *testing.T) {
	// Prepare to assert multiple times.
	assert := assert.New(t)
	raw, err := acigo.NewACIFaultSeverityRaw("info")
	assert.Nil(err)
	assert.Equal("info", fmt.Sprintf("%s", raw))
}
