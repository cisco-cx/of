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

package v1

// LabelName represents a key in a LabelMap.
//
// LabelName is based on LabelName in `github.com/prometheus/common/model`.
type LabelName string

// LabelNameValidator is an interface that represents the ability to
// return a non-nil error if a given LabelName is invalid.
type LabelNameValidator interface {
	Validate() error
}
