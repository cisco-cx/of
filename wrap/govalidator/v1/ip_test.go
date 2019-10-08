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
// The MIT License (MIT)
//
// Copyright (c) 2014 Alex Saskevich
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package v1_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	of "github.com/cisco-cx/of/lib/v1"
	govalidator "github.com/cisco-cx/of/wrap/govalidator/v1"
)

// IPv4

// Confirm that govalidator.IP implements the of.IPv4Validator interface.
func TestIP_InterfaceIPv4Validator(t *testing.T) {
	var _ of.IPv4Validator = &govalidator.IP{}
	assert.Nil(t, nil) // If we get this far, the test passed.
}

// Check a simple negative path of govalidator.IsIPv4.
func TestIP_IPv4SimpleNegative(t *testing.T) {
	assert := assert.New(t) // Prepare to assert multiple times.
	ip, err := govalidator.NewIP("2001:db8::1")
	assert.Nil(err)
	assert.False(ip.IsIPv4())
}

// Check a simple positive path of govalidator.IsIPv4.
func TestIP_IPv4SimplePositive(t *testing.T) {
	assert := assert.New(t) // Prepare to assert multiple times.
	ip, err := govalidator.NewIP("192.168.222.222")
	assert.Nil(err)
	assert.True(ip.IsIPv4())
}

// IPv6

// Confirm that govalidator.IP implements the of.IPv6Validator interface.
func TestIP_InterfaceIPv6Validator(t *testing.T) {
	var _ of.IPv6Validator = &govalidator.IP{}
	assert.Nil(t, nil) // If we get this far, the test passed.
}

// Check a simple negative path of govalidator.IsIPv6.
func TestIP_IPv6SimpleNegative(t *testing.T) {
	assert := assert.New(t) // Prepare to assert multiple times.
	ip, err := govalidator.NewIP("192.168.222.222")
	assert.Nil(err)
	assert.False(ip.IsIPv6())
}

// Check a simple positive path of govalidator.IsIPv6.
func TestIP_IPv6SimplePositive(t *testing.T) {
	assert := assert.New(t) // Prepare to assert multiple times.
	ip, err := govalidator.NewIP("2001:db8::1")
	assert.Nil(err)
	assert.True(ip.IsIPv6())
}

// Generic

// Simple test of IP's implementation of the fmt.Stringer interface.
func TestGovalidatorIP_StringerSimple(t *testing.T) {
	// Prepare to assert multiple times.
	assert := assert.New(t)
	ip, err := govalidator.NewIP("2001:db8::1")
	assert.Nil(err)
	assert.Equal("2001:db8::1", fmt.Sprintf("%s", ip))
}
