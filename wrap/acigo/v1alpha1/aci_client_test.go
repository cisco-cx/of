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

package v1alpha1_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	of "github.com/cisco-cx/of/lib/v1alpha1"
	acigo "github.com/cisco-cx/of/wrap/acigo/v1alpha1"
)

// Confirm that acigo.ACIClient implements the of.ACIClient interface.
func TestACIClient_Interface(t *testing.T) {
	var _ of.ACIClient = &acigo.ACIClient{}
	assert.Nil(t, nil) // If we get this far, the test passed.
}

// Confirm simple acigo.NewACIClient() functionality.
func TestACIClient_Simple(t *testing.T) {
	// Get a new client.
	_, err := acigo.NewACIClient(of.ACIClientConfig{
		Hosts: []string{"host1", "host2", "host3"},
		User:  "user",
		Pass:  "pass",
	})
	assert.Nil(t, err)
}

// Simple ACIClient.Faults() functionality test.
func TestACIClient_FaultsSimple(t *testing.T) {
	// Prepare to assert multiple times.
	assert := assert.New(t)
	// Get a new client and confirm nil err.
	client, err := acigo.NewACIClient(of.ACIClientConfig{
		Hosts: []string{"host1", "host2", "host3"},
		User:  "user",
		Pass:  "pass",
	})
	// Confirm nil err for new client.
	assert.Nil(err)
	// Get ACI faults without type inference.
	var faults []of.Map
	faults, err = client.Faults()
	// Confirm nil err for getting ACI faults.
	assert.Nil(err)
	// Confirm faults data type (again).
	assert.Equal([]of.Map{}, faults)
}
