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

package v1alpha1

// LabelMap is a Domain Type that represents a map of labels used to label a
// resource. It is based on `prometheus/common.LabelSet`
type LabelMap map[LabelName]LabelValue

// LabelMapCopier is an interface that represents the ability to return a
// copy of an existing LabelMap.
type LabelMapCopier interface {
	// Copy returns a copy of an existing LabelMap.
	Copy() LabelMap
}

// LabelMapEqualer is an interface that represents the ability to return
// true only if two LabelMap structs hold the same data.
type LabelMapEqualer interface {
	// Equal returns true when two LabelMap structs hold the same data.
	Equal(other LabelMap) bool
}

// LabelMapValidator is an interface that represents the ability to
// return a non-nil error if a given LabelMap is invalid.
type LabelMapValidator interface {
	// Validate returns a non-nil error if a given LabelMap is found to be
	// invalid.
	Validate() error
}

// LabelMapMerger is an interface that represents the ability to return
// a new merged LabelMap, given two given LabelMap structs.
type LabelMapMerger interface {
	// Merge returns a new merged LabelMap, given two LabelMap structs.
	Merge(other LabelMap) LabelMap
}

// LabelMapStringer is an interface that embeds the `fmt.Stringer`
// interface.
type LabelMapStringer interface {
	// String implements the `fmt.Stringer` interface.
	String() string
}
