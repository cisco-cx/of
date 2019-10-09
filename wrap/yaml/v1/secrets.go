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
	"io"

	"gopkg.in/yaml.v2"
	of "github.com/cisco-cx/of/pkg/v1"
)

type Secrets of.Secrets

func (s *Secrets) Decode(r io.Reader) error {
	return yaml.NewDecoder(r).Decode(s)
}

func (s *Secrets) Encode(w io.Writer) error {
	return yaml.NewEncoder(w).Encode(s)
}
