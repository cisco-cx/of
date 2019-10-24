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

package v1

import (
	"bytes"
	"io"
	"io/ioutil"

	"gopkg.in/yaml.v2"
	snmp_config "github.com/cisco-cx/of/pkg/v2/snmp"
)

type Configs snmp_config.V2Config

// Implements snmp v2 config Decoder.
func (a *Configs) Decode(r io.Reader) error {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(data, a)
}

// Implements snmp v2 config Encoder.
func (a *Configs) Encode(w io.Writer) error {
	data, err := yaml.Marshal(a)
	if err != nil {
		return err
	}
	r := bytes.NewReader(data)
	if _, err = io.Copy(w, r); err != nil {
		return err
	}

	return nil
}
