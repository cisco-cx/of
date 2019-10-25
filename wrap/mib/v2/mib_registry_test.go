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
package v2_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	of "github.com/cisco-cx/of/pkg/v2"
	mib "github.com/cisco-cx/of/wrap/mib/v2"
)

func TestHealthChecker_Interface(t *testing.T) {
	wrapMib := mib.New()
	var _ of.MibRegistry = wrapMib
}

var testMib of.MibRegistry

func init() {
	testMib = mib.New()
	m := map[string]of.Mib{
		"1":      of.Mib{Name: "1st-Name", Description: "I am 1", Units: "bytes"},
		"1.2":    of.Mib{Name: "2nd-Name", Description: "I am 2", Units: "bytes"},
		"1.2.3":  of.Mib{Name: "3rd-Name", Description: "I am 3", Units: "bytes"},
		".1.2":   of.Mib{Name: "5th-Name", Description: "I am 5", Units: "bytes"},
		".1.2.6": of.Mib{Name: "6th-Name", Description: "I am 6", Units: "bytes"},
		".1.2.8": of.Mib{Name: "7th-Name", Description: "I am 7", Units: "bytes"},
	}
	mibs := map[string]of.Mib{
		".1.3.6.1.2.1.1.3.0": of.Mib{
			Name: "oid1",
		},
		".1.3.6.1.6.3.1.1.4.1.0": of.Mib{
			Name: "oid2",
		},
		".1.3.6.1.4.1.8164.2.44": of.Mib{
			Name: "oid3",
		},
		".1.3.6.1.4.1.8164.2.45": of.Mib{
			Name: "oid4",
		},
		".1.3.6.1.4.1.65000.1.1.1.1.1": of.Mib{
			Name: "oid5",
		},
	}
	testMib.Load(m)
	testMib.Load(mibs)
}

func TestMibRegistry_Load(t *testing.T) {
	testLoad := mib.New()
	m := map[string]of.Mib{
		"1":     of.Mib{Name: "1st-Name", Description: "I am 1", Units: "bytes"},
		"1.2":   of.Mib{Name: "2nd-Name", Description: "I am 2", Units: "bytes"},
		"1.2.3": of.Mib{Name: "3rd-Name", Description: "I am 3", Units: "bytes"},
	}
	error := testLoad.Load(m)
	// load successfully
	require.NoError(t, error)

	m["3.4"] = of.Mib{Description: "my desc"}
	error = testLoad.Load(m)
	// name must be present
	assert.EqualError(
		t,
		error,
		"Name can't be empty: '{Name: Description:my desc Units:}'",
	)
}

func TestMibRegistry_String_lastNodeMissing(t *testing.T) {
	value := testMib.String("1.2.7")
	assert.Equal(t, "1st-Name.2nd-Name.7", value)
}

func TestMibRegistry_String_fullMatch(t *testing.T) {
	value := testMib.String("1.2")
	assert.Equal(t, "1st-Name.2nd-Name", value)
}

func TestMibRegistry_String_singleNode(t *testing.T) {
	value := testMib.String("single node")
	assert.Equal(t, "single node", value)
}

func TestMibRegistry_String_empty(t *testing.T) {
	value := testMib.String("")
	assert.Equal(t, "", value)
}

func TestMibRegistry_String_oidNotFound(t *testing.T) {
	value := testMib.String("2.3.repeat after me")
	assert.Equal(t, "2.3.repeat after me", value)
}

func TestMibRegistry_String_nodeInTheMiddleFound(t *testing.T) {
	value := testMib.String(".1.2.4")
	assert.Equal(t, ".1.5th-Name.4", value)
}

func TestMibRegistry_String_startWithDot(t *testing.T) {
	value := testMib.String(".1.3.6.1.2.1.1.3.0")
	assert.Equal(t, ".1.3.6.1.2.1.1.3.oid1", value)
}

func TestMibRegistry_ShortString_fullMatch(t *testing.T) {
	value := testMib.ShortString("1.2")
	assert.Equal(t, "2nd-Name", value)
}

func TestMibRegistry_ShortString_fullMatchStartsWithDot(t *testing.T) {
	value := testMib.ShortString(".1.2.8")
	assert.Equal(t, "7th-Name", value)
}

func TestMibRegistry_ShortString_notFound(t *testing.T) {
	value := testMib.ShortString("not found")
	assert.Equal(t, "", value)
}

func TestMibRegistry_ShortString_rootNodesLoaded(t *testing.T) {
	value := testMib.ShortString("1.2.3.not loaded but parents")
	assert.Equal(t, "", value)
}
