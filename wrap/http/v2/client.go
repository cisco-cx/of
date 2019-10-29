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
	"io"
	"net/http"
)

// Wrapping http.Client
func NewClient() *http.Client {
	return &http.Client{}
}

// Wrapping http.Requests
func NewRequest(method, url string, body io.Reader) (*http.Request, error) {
	return http.NewRequest(method, url, body)
}
