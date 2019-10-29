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
	"bufio"
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	of "github.com/cisco-cx/of/pkg/v1"
	cache "github.com/cisco-cx/of/wrap/cache/v1"
)

func TestCacher_Interface(t *testing.T) {
	var c cache.Cacher
	var _ of.Cacher = &c
}

func TestCacher_Read(t *testing.T) {
	var c cache.Cacher
	r := strings.NewReader(`{"name": "foo"}`)
	m := map[string]string{}

	err := c.Read(r, &m)
	require.NoError(t, err)
	require.Equal(t, "foo", m["name"])

	r.Reset(`{"name": "bar"}`)
	var mib of.MIB
	err = c.Read(r, &mib)
	require.NoError(t, err)
	require.Equal(t, "bar", mib.Name)
}

func TestCacher_Write(t *testing.T) {
	var c cache.Cacher
	m := map[string]string{
		"name": "foo",
	}

	var buf bytes.Buffer
	w := bufio.NewWriter(&buf)
	err := c.Write(w, &m)
	require.NoError(t, err)
	w.Flush()
	expected := bytes.NewBufferString(`{"name":"foo"}
`)
	require.Equal(t, expected, &buf)
}
