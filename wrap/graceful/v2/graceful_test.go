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

package v2_test

import (
	"io/ioutil"
	"net/http"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"gotest.tools/assert"
	of "github.com/cisco-cx/of/pkg/v2"
	graceful "github.com/cisco-cx/of/wrap/graceful/v2"
)

// Represents a net/http.Handler.
type testServer struct {
	timeout time.Duration
}

// Implementing  net/http.Handler interface.
func (s *testServer) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	time.Sleep(s.timeout)
	rw.Write([]byte("hi"))
}

// Enforcing interface implementation.
func TestInterface(t *testing.T) {
	var _ of.Graceful = &graceful.Graceful{}
}

// Test gracefully starting a server.
func TestGraceStart(t *testing.T) {
	server := &http.Server{
		Addr:    "localhost:54951",
		Handler: &testServer{timeout: time.Second * 0},
	}

	g := graceful.New(server)
	go func() {
		require.NoError(t, g.Start())
	}()
	time.Sleep(time.Second)

	res, err := http.Get("http://localhost:54951/")

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	all, err := ioutil.ReadAll(res.Body)
	require.NoError(t, err)
	assert.Equal(t, "hi", string(all))

	require.NoError(t, g.Stop())
}

// Test gracefully shuting down a server.
func TestGraceStop(t *testing.T) {
	server := &http.Server{
		Addr:    "localhost:54952",
		Handler: &testServer{timeout: time.Second * 10},
	}

	g := graceful.New(server)
	go func() {
		require.NoError(t, g.Start())
	}()
	time.Sleep(time.Second)

	require.NoError(t, g.Stop())

	_, err := http.Get("http://localhost:54952/")
	require.Error(t, err)
}

// Test with SIGINT a server.
func TestKill(t *testing.T) {
	server := &http.Server{
		Addr:    "localhost:54953",
		Handler: &testServer{timeout: time.Second * 10},
	}

	g := graceful.New(server)
	go func() {
		require.NoError(t, g.Start())
	}()
	time.Sleep(time.Second)

	_, err := http.Get("http://localhost:54953/")
	require.Error(t, err)
	syscall.Kill(syscall.Getpid(), syscall.SIGINT)
}
