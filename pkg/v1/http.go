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
	"net/http"
	"time"
)

type ResponseWriter http.ResponseWriter
type Request *http.Request

type HTTPConfig struct {
	ListenAddress string
	ReadTimeout   time.Duration
	WriteTimeout  time.Duration
}

type Server struct {
	Srv *http.Server
	Mux *http.ServeMux
	G   Graceful
}

// Represents HTTP handler
type Handler interface {
	ServeHTTP(ResponseWriter, Request)
}

// Represents HTTP server components
type Serve interface {
	ListenAndServe() error
	Shutdown() error
	Handle(string, Handler)
	HandleFunc(string, func(ResponseWriter, Request))
}

// Represents HTTP client
type Client http.Client
