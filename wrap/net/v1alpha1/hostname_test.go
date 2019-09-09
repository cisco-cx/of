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
// Copyright (c) 2009 The Go Authors. All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are
// met:
//
//    * Redistributions of source code must retain the above copyright
// notice, this list of conditions and the following disclaimer.
//    * Redistributions in binary form must reproduce the above
// copyright notice, this list of conditions and the following disclaimer
// in the documentation and/or other materials provided with the
// distribution.
//    * Neither the name of Google Inc. nor the names of its
// contributors may be used to endorse or promote products derived from
// this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
// OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
// LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
// This work incorporates works covered by the following notices:

package v1alpha1_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	of "github.com/cisco-cx/of/lib/v1alpha1"
	net "github.com/cisco-cx/of/wrap/net/v1alpha1"
)

// Confirm that Hostname implements the of.IPv4Finder interface.
func TestHostname_InterfaceIPv4Finder(t *testing.T) {
	var _ of.IPv4Finder = &net.Hostname{}
	assert.Nil(t, nil) // If we get this far, the test passed.
}

// Confirm that Hostname implements the of.IPv6Finder interface.
func TestHostname_InterfaceIPv6Finder(t *testing.T) {
	var _ of.IPv6Finder = &net.Hostname{}
	assert.Nil(t, nil) // If we get this far, the test passed.
}

// Simple test of Hostname's implementation of the fmt.Stringer interface.
func TestHostname_StringerSimple(t *testing.T) {
	// Prepare to assert multiple times.
	assert := assert.New(t)
	h, err := net.NewHostname("hello.example.org")
	assert.Nil(err)
	assert.Equal("hello.example.org", fmt.Sprintf("%s", h))
}
