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
	"context"
	"net/http"
	"time"

	"github.com/ory/graceful"
)

// Holds *http.Server that needs to be stopped graceful.
type Graceful struct {
	server *http.Server
}

// Initialize Graceful
func New(srv *http.Server) *Graceful {
	g := &Graceful{
		server: graceful.WithDefaults(srv),
	}
	return g
}

// Starts a http.Server that shuts down on SIGINT or SIGTERM.
func (g *Graceful) Start() error {
	return graceful.Graceful(g.start, g.stop)
}

// Callback executed by graceful to start the server.
func (g *Graceful) start() error {
	if len(g.server.TLSConfig.Certificates) == 0 {
		return g.server.ListenAndServe()
	}
	return g.server.ListenAndServeTLS("", "")
}

// Callback to be called on SIGINT or SIGTERM.
func (g *Graceful) stop(ctx context.Context) error {
	return g.server.Shutdown(ctx)
}

// Shuts down the server.
func (g *Graceful) Stop() error {
	timer, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return g.server.Shutdown(timer)
}
