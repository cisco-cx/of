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
// Copyright (c) 2013 Mitchell Hashimoto
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, Subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or Substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package v1alpha1_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	mapstructure "github.com/cisco-cx/of/wrap/mapstructure/v1alpha1"
)

type TestPerson struct {
	Name   string
	Age    int
	Emails []string
	Extra  map[string]string
}

// Ensure a simple Map is decoded as expected.
//
// TestMap_DecodeSimple is based on:
// https://godoc.org/github.com/mitchellh/mapstructure#ex-Decode
func TestMap_DecodeSimple(t *testing.T) {
	// Prepare to assert multiple times.
	assert := assert.New(t)

	// Create a simple Map without type-inference.
	input := map[string]interface{}{
		"name":   "Mitchell",
		"age":    91,
		"emails": []string{"one", "two", "three"},
		"extra": map[string]string{
			"twitter": "mitchellh",
		},
	}
	m := mapstructure.NewMap(input)

	// Declare expected outcome.
	var expect TestPerson = TestPerson{
		Name:   "Mitchell",
		Age:    91,
		Emails: []string{"one", "two", "three"},
		Extra: map[string]string{
			"twitter": "mitchellh",
		},
	}

	// Confirm the decode went as planned.
	var result TestPerson
	err := m.DecodeMap(&result)
	assert.Equal(err, nil, "mapstructure.Decode() returned non-nil error")
	assert.Equal(expect, result, "Did not obtain expected result.")
}
