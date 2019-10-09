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
	"testing"

	"github.com/stretchr/testify/assert"

	of "github.com/cisco-cx/of/pkg/v1"
	acigo "github.com/cisco-cx/of/wrap/acigo/v1"
)

// Confirm that acigo.ACIClient implements the of.ACIClient interface.
func TestACIClient_Interface(t *testing.T) {
	var _ of.ACIClient = &acigo.ACIClient{}
	assert.Nil(t, nil) // If we get this far, the test passed.
}
