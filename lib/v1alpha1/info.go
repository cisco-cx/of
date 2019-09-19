// Copyright 2019 Cisco Systems, Inc.
//
// This work incorporates works covered by the following notices:
// ---
// Copyright 2016 The Prometheus Authors
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
// ---
// The MIT License (MIT)
//
// Copyright (c) 2017 Middlemost Systems, LLC
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

package v1alpha1

// InfoService is a domain interface that represents a service for exposing build and version information.
type InfoService interface {
	// BuildInfo returns build information as a string.
	BuildInfo() string

	// Metadata returns metadata as a string.
	Metadata() string

	// String returns metadata, build and version information as a string.
	String() string

	// VersionInfo returns version information as a string.
	VersionInfo() string
}

// InformerCollector is an optional extension to InfoService that represents the ability to expose metadata, build and
// version information as a Prometheus metrics collector.
type InfoCollector interface {
	// Register registers a collector with Prometheus that exports metrics about InfoService.
	Register() error

	// Unregister unregisters the collector with Prometheus. Normally, this would be called in a defer statement.
	Unregister() bool
}
