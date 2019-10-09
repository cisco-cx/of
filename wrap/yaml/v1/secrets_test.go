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
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	of "github.com/cisco-cx/of/pkg/v1"
	yaml "github.com/cisco-cx/of/wrap/yaml/v1"
)

// Enforce interface implementation.
func TestSecretsInterface(t *testing.T) {
	var _ of.Decoder = &yaml.Secrets{}
	var _ of.Encoder = &yaml.Secrets{}
}

// Ensure yaml decodes Secrets
func TestSecretsDecoder(t *testing.T) {

	r := strings.NewReader(`
apic:
  cluster:
    name: lab-aci`)

	expected := of.Secrets{
		APIC: of.SecretsConfigAPIC{
			Cluster: of.SecretsConfigAPICCluster{
				Name: "lab-aci",
			},
		},
	}

	cfg := yaml.Secrets{}
	cfg.Decode(r)
	require.EqualValues(t, expected, cfg)
}

// Ensure yaml encodes Secrets
func TestSecretsEncoder(t *testing.T) {

	expected := `apic:
  cluster:
    name: lab-aci`

	cfg := yaml.Secrets{
		APIC: of.SecretsConfigAPIC{
			Cluster: of.SecretsConfigAPICCluster{
				Name: "lab-aci",
			},
		},
	}

	buf := bytes.NewBuffer(nil)
	cfg.Encode(buf)
	require.EqualValues(t, expected, strings.Trim(string(buf.Bytes()), "\n"))
}
