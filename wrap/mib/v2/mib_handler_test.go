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
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	of "github.com/cisco-cx/of/pkg/v2"
	mib "github.com/cisco-cx/of/wrap/mib/v2"
)

func TestMIBHandler_LoadJSONFromFile_invalidPath(t *testing.T) {
	testMIBHandler := &mib.MIBHandler{
		MapMIB: make(map[string]of.MIB),
	}
	error := testMIBHandler.LoadJSONFromFile("assets/doesntexist.json")
	require.Error(t, error)
}

func TestMIBHandler_LoadJSONFromFile_invalidFormat(t *testing.T) {
	testMIBHandler := &mib.MIBHandler{
		MapMIB: make(map[string]of.MIB),
	}
	error := testMIBHandler.LoadJSONFromFile("assets/invalid-format.json")
	require.Error(t, error)
}

func TestMIBHandler_LoadJSONFromFile_ok(t *testing.T) {
	testMIBHandler := &mib.MIBHandler{
		MapMIB: make(map[string]of.MIB),
	}
	error := testMIBHandler.LoadJSONFromFile("assets/json/SPIDCOM-MIB.json")
	require.NoError(t, error)
	require.Equal(t, "spidcom", testMIBHandler.MapMIB["1.3.6.1.4.1.22764"].Name)
}

func TestMIBHandler_LoadJSONFromDir_dirNotFound(t *testing.T) {
	testMIBHandler := &mib.MIBHandler{
		MapMIB: make(map[string]of.MIB),
	}
	error := testMIBHandler.LoadJSONFromDir("assets/not-found")
	require.Error(t, error)
}

func TestMIBHandler_LoadJSONFromDir_notDir(t *testing.T) {
	testMIBHandler := &mib.MIBHandler{
		MapMIB: make(map[string]of.MIB),
	}
	error := testMIBHandler.LoadJSONFromDir("assets/invalid-format.json")
	require.Error(t, error)
	assert.Contains(t, error, "not a directory")
}

func TestMIBHandler_LoadJSONFromDir_ok(t *testing.T) {
	testMIBHandler := &mib.MIBHandler{
		MapMIB: make(map[string]of.MIB),
	}
	error := testMIBHandler.LoadJSONFromDir("assets/json")
	require.NoError(t, error)
	require.Equal(t, "spidcom", testMIBHandler.MapMIB["1.3.6.1.4.1.22764"].Name)
}

func TestMIBHandler_WriteCacheToFile_invalidPath(t *testing.T) {
	testMIBHandler := &mib.MIBHandler{
		MapMIB: make(map[string]of.MIB),
	}
	error := testMIBHandler.LoadJSONFromDir("assets/json")
	require.NoError(t, error)

	error = testMIBHandler.WriteCacheToFile("/no/such/file/or/directory")
	require.Error(t, error)
}

func TestMIBHandler_WriteCacheToFile_ok(t *testing.T) {
	testMIBHandler := &mib.MIBHandler{
		MapMIB: make(map[string]of.MIB),
	}
	error := testMIBHandler.LoadJSONFromDir("assets/json")
	require.NoError(t, error)

	file, err := ioutil.TempFile("/tmp", "mibHandler")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(file.Name())

	error = testMIBHandler.WriteCacheToFile(file.Name())
	require.NoError(t, error)

	data, err := ioutil.ReadFile(file.Name())
	require.NoError(t, err)

	assert.Contains(t, string(data), "spidcom")
}

func TestMIBHandler_LoadCacheFromFile_ok(t *testing.T) {
	readerMIB := &mib.MIBHandler{
		MapMIB: make(map[string]of.MIB),
	}
	error := readerMIB.LoadCacheFromFile("assets/mib.cache")
	require.NoError(t, error)
	require.Equal(t, "spidcom", readerMIB.MapMIB["1.3.6.1.4.1.22764"].Name)
}

func TestMIBHandler_LoadCacheFromFile_invalidPath(t *testing.T) {
	readerMIB := &mib.MIBHandler{
		MapMIB: make(map[string]of.MIB),
	}
	error := readerMIB.LoadCacheFromFile("/no/such/path/or/directory")
	require.Error(t, error)
}
