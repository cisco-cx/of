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
	"io"
	"sort"
	"testing"

	"github.com/stretchr/testify/require"
	of "github.com/cisco-cx/of/pkg/v2"
	concatenator "github.com/cisco-cx/of/wrap/concatenator/v2"
	yaml "github.com/cisco-cx/of/wrap/yaml/v2"
)

// Enforce interface implementation.
func TestConcatenateInterface(t *testing.T) {
	var _ of.Concatenate = &concatenator.Files{}
}

// Test Concat for different paths.
func TestConcat(t *testing.T) {
	// Expects yaml decode to succeed
	verifyConcat(t, "testdata/", "", verifyYaml)

	// Expects yaml decode to succeed
	verifyConcat(t, "testdata", "", verifyYaml)

	// Expects yaml decode to succeed
	verifyConcat(t, "testdata", "yaml", verifyYaml)

	// Expects no errors, when no files are found for given extension.
	verifyConcat(t, "testdata", "yaml1", requireNoError)

	// Expects errors, when path is not a directory.
	verifyConcat(t, "somerandompath", "yaml1", requireError)
}

// Concat files in given path and call testFunc to verify expected result.
func verifyConcat(t *testing.T, path string, ext string, testFunc func(t *testing.T, r io.Reader, err error)) {

	files := concatenator.Files{Path: path, Ext: ext}
	r, err := files.Concat()
	testFunc(t, r, err)
}

// Decode using yaml to verify content sanity.
func verifyYaml(t *testing.T, r io.Reader, err error) {

	require.NoError(t, err)

	cfgs := yaml.Configs{}
	err = cfgs.Decode(r)
	require.NoError(t, err)
	expectedNames := []string{"oki_ntp", "cpnr_epm", "cpnr_system", "epc", "esc", "nso"}
	names := make([]string, 0)
	for cfgName, _ := range cfgs {
		names = append(names, cfgName)
	}
	sort.Strings(names)
	sort.Strings(expectedNames)
	require.Equal(t, expectedNames, names)
}

// Expects err to be none nil.
func requireError(t *testing.T, r io.Reader, err error) {
	require.Error(t, err)
}

// Expects err to be nil.
func requireNoError(t *testing.T, r io.Reader, err error) {
	require.NoError(t, err)
}
