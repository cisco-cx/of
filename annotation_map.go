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

package of

// AnnotationMap is a Domain Type that represents a LabelMap used to annotate
// rather than label a resource.
type AnnotationMap LabelMap

// AnnotationMapCopier is an interface that represents the ability to return a
// copy of an existing AnnotationMap.
type AnnotationMapCopier interface {
	// Copy returns a copy of an existing AnnotationMap.
	Copy() AnnotationMap
}

// AnnotationMapEqualer is an interface that represents the ability to return
// true only if two AnnotationMap structs hold the same data.
type AnnotationMapEqualer interface {
	// Equal returns true when two AnnotationMap structs hold the same data.
	Equal(other AnnotationMap) bool
}

// AnnotationMapValidator is an interface that represents the ability to
// return a non-nil error if a given AnnotationMap is invalid.
type AnnotationMapValidator interface {
	// Validate returns a non-nil error if a given AnnotationMap is found to be
	// invalid.
	Validate() error
}

// AnnotationMapMerger is an interface that represents the ability to return
// a new merged AnnotationMap, given two given AnnotationMap structs.
type AnnotationMapMerger interface {
	// Merge returns a new merged AnnotationMap, given two AnnotationMap structs.
	Merge(other AnnotationMap) AnnotationMap
}

// AnnotationMapStringer is an interface that embeds the `fmt.Stringer`
// interface.
type AnnotationMapStringer interface {
	// String implements the `fmt.Stringer` interface.
	String() string
}
