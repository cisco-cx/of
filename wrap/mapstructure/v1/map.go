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

package v1

import (
	"github.com/mitchellh/mapstructure"

	of "github.com/cisco-cx/of/lib/v1"
)

// Map represents an arbitrary map[string]interface{} data that will be decoded
// into a native Go structure.
//
// Map is based on Map in:
// "github.com/cisco-cx/of/lib/v1"
type Map struct {
	ofMap of.Map
}

// NewMap returns a new instance of Map.
func NewMap(input map[string]interface{}) Map {
	return Map{
		ofMap: input,
	}
}

// DecodeMap decodes a raw interface into structured data.
//
// DecodeMap implements MapDecoder in:
// "github.com/cisco-cx/of/lib/v1"
//
// DecodeMap is based on Decode in: "github.com/mitchellh/mapstructure"
func (m Map) DecodeMap(output interface{}) error {
	// Convert m to an interface{} var.
	var input interface{} = m.ofMap
	// Call external package to decode input into output.
	err := mapstructure.Decode(input, &output)
	return err
}
