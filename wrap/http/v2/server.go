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

	of "github.com/cisco-cx/of/pkg/v2"
	graceful "github.com/cisco-cx/of/wrap/graceful/v2"
)

// Represents server components.
type Server of.Server

// Initialize a server.
func NewServer(config *of.HTTPConfig) *Server {

	srv := &http.Server{
		Addr:         config.ListenAddress,
		ReadTimeout:  config.ReadTimeout,
		WriteTimeout: config.WriteTimeout,
	}

	m := http.NewServeMux()
	srv.Handler = m
	return &Server{
		Srv: srv,
		Mux: m,
		G:   graceful.New(srv)}
}

// Start a graceful server.
func (s *Server) ListenAndServe() error {
	return s.G.Start()
}

// Shutdown http server.
func (s *Server) Shutdown() error {
	return s.G.Stop()
}

// Handle registers the handler for the given pattern.
func (s *Server) Handle(pattern string, h of.Handler) {
	newHandler := &handlerOveride{h.ServeHTTP}
	s.Mux.Handle(pattern, newHandler)
}

// HandleFunc registers the handler function for the given pattern.
func (s *Server) HandleFunc(pattern string, h func(of.ResponseWriter, of.Request)) {
	newHandler := func(rw http.ResponseWriter, r *http.Request) {
		h(rw, r)
	}
	s.Mux.HandleFunc(pattern, newHandler)
}
