// Copyright 2019 Cisco Systems, Inc.
//
// This work incorporates works covered by the following notice:
//
// Copyright © 2015-2018 Aeneas Rekkas <aeneas+oss@aeneas.io>

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

package v2

// Writer is a helper to write arbitrary data to a ResponseWriter
type Writer interface {
	// Write a response object to the ResponseWriter with status code 200.
	Write(w ResponseWriter, r Request, e interface{})

	// WriteCode writes a response object to the ResponseWriter and sets a response code.
	WriteCode(w ResponseWriter, r Request, code int, e interface{})

	// WriteCreated writes a response object to the ResponseWriter with status code 201 and
	// the Location header set to location.
	WriteCreated(w ResponseWriter, r Request, location string, e interface{})

	// WriteError writes an error to ResponseWriter and tries to extract the error's status code by
	// asserting statusCodeCarrier. If the error does not implement statusCodeCarrier, the status code
	// is set to 500.
	WriteError(w ResponseWriter, r Request, err interface{})

	// WriteErrorCode writes an error to ResponseWriter and forces an error code.
	WriteErrorCode(w ResponseWriter, r Request, code int, err interface{})
}
