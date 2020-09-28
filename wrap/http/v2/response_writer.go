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

package v2

import (
	"net/http"
)

// Wrapping net/http ResponseWriter
type ResponseWriter struct {
	hrw        http.ResponseWriter
	statusCode int
}

// Initiate of.ResponseWriter
func NewResponseWriter(hrw http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{hrw, 200}
}

// Wrapping net/http ResponseWriter method
func (rw *ResponseWriter) Header() http.Header {
	return rw.hrw.Header()
}

// Wrapping net/http ResponseWriter method
func (rw *ResponseWriter) Write(b []byte) (int, error) {
	return rw.hrw.Write(b)
}

// Wrapping net/http ResponseWriter method and storing status code
func (rw *ResponseWriter) WriteHeader(i int) {
	rw.statusCode = i
	rw.hrw.WriteHeader(i)
}

// Return statusCode saved by `WriteHeader` method
func (rw *ResponseWriter) StatusCode() int {
	return rw.statusCode
}
