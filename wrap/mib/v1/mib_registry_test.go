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
package v1_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	of "github.com/cisco-cx/of/pkg/v1"
	mib "github.com/cisco-cx/of/wrap/mib/v1"
)

func TestHealthChecker_Interface(t *testing.T) {
	wrapMib := mib.New()
	var _ of.MibRegistry = wrapMib
}

var testMib of.MibRegistry

func init() {
	testMib = mib.New()
	m := map[string]of.Mib{
		"1":     of.Mib{Name: "1st-Name", Description: "I am 1", Units: "bytes"},
		"1.2":   of.Mib{Name: "2st-Name", Description: "I am 2", Units: "bytes"},
		"1.2.3": of.Mib{Name: "3st-Name", Description: "I am 3", Units: "bytes"},
	}
	testMib.Load(m)
}

func TestMibRegistry_Load(t *testing.T) {
	testMib = mib.New()
	m := map[string]of.Mib{
		"1":     of.Mib{Name: "1st-Name", Description: "I am 1", Units: "bytes"},
		"1.2":   of.Mib{Name: "2nd-Name", Description: "I am 2", Units: "bytes"},
		"1.2.3": of.Mib{Name: "3rd-Name", Description: "I am 3", Units: "bytes"},
	}
	error := testMib.Load(m)
	require.NoError(t, error)

	m["3.4"] = of.Mib{Description: "my desc"}
	error = testMib.Load(m)
	assert.EqualError(
		t,
		error,
		"Name can't be empty: '{Name: Description:my desc Units:}'",
	)
}

func TestMibRegistry_String(t *testing.T) {
	value := testMib.String("1.2")
	assert.Equal(t, "1st-Name.2nd-Name", value)

	value = testMib.String("invalid name")
	assert.Equal(t, "", value)

	value = testMib.String("2.3.not loaded")
	assert.Equal(t, "", value)

	value = testMib.String("1.2.3.not loaded but parents")
	assert.Equal(t, "1st-Name.2nd-Name.3rd-Name.not loaded but parents", value)

	value = testMib.String("1.3.4")
	assert.Equal(t, "1st-Name.3.4", value)
}

func TestMibRegistry_ShortString(t *testing.T) {
	value := testMib.ShortString("1.2")
	assert.Equal(t, "2nd-Name", value)

	value = testMib.ShortString("invalid name")
	assert.Equal(t, "", value)

	value = testMib.ShortString("2.3.not loaded")
	assert.Equal(t, "", value)

	value = testMib.ShortString("1.2.3.not loaded but parents")
	assert.Equal(t, "", value)

	value = testMib.ShortString("1.3.4")
	assert.Equal(t, "", value)
}
