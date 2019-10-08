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

// APIVersion represents a named version of an API (e.g. "v1", "v1").
type APIVersion string

// APIVersionValidator is an interface that represents the ability to
// return a non-nil error if a given APIVersion is invalid.
type APIVersionValidator interface {
	// Validate returns a non-nil error if a given APIVersion is found to be
	// invalid.
	Validate() error
}

// Validate returns a non-nil error if a given APIVersion is found to be
// invalid.
// func (v APIVersion) Validate() error
