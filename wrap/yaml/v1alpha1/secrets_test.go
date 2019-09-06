// Copyright 2019 Cisco Systems, Inc.
//
// This work incorporates works covered by the following notice:
//
// Copyright 2018 The Prometheus Authors
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

package v1alpha1_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/cisco-cx/of/lib/v1alpha1"
	yaml "github.com/cisco-cx/of/wrap/yaml/v1alpha1"
)

// Ensure yaml decodes Secrets
func TestSecretsDecoder(t *testing.T) {

	r := strings.NewReader(`
apic:
  cluster:
    name: lab-aci`)

	expected := v1alpha1.Secrets{
		APIC: v1alpha1.SecretsConfigAPIC{
			Cluster: v1alpha1.SecretsConfigAPICCluster{
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
		APIC: v1alpha1.SecretsConfigAPIC{
			Cluster: v1alpha1.SecretsConfigAPICCluster{
				Name: "lab-aci",
			},
		},
	}

	buf := bytes.NewBuffer(nil)
	cfg.Encode(buf)
	require.EqualValues(t, expected, strings.Trim(string(buf.Bytes()), "\n"))
}
