// Copyright 2019 Cisco Systems, Inc.
//
// This work incorporates works covered by the following notice:
//
//The MIT License (MIT)

//Copyright (c) 2017 InVision

//Permission is hereby granted, free of charge, to any person obtaining a copy
//of this software and associated documentation files (the "Software"), to deal
//in the Software without restriction, including without limitation the rights
//to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
//copies of the Software, and to permit persons to whom the Software is
//furnished to do so, subject to the following conditions:

//The above copyright notice and this permission notice shall be included in all
//copies or substantial portions of the Software.

//THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
//IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
//FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
//AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
//LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
//OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
//SOFTWARE.

package v2_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	of "github.com/cisco-cx/of/pkg/v1"
	health "github.com/cisco-cx/of/wrap/health/v1"
)

func TestHealthChecker_Interface(t *testing.T) {
	wrapHealth := health.New()
	var _ of.HealthChecker = wrapHealth
}

func TestHealthChecker_Start(t *testing.T) {
	wrapHealth := health.New()
	err := wrapHealth.Start()
	require.NoError(t, err)
	wrapHealth.Stop()
}

func TestHealthChecker_Stop(t *testing.T) {
	wrapHealth := health.New()
	err := wrapHealth.AddURL("name", "http://localhost/somerandomurl", 2*time.Second)
	require.NoError(t, err)

	err = wrapHealth.Stop()
	assert.EqualError(
		t,
		err,
		"Healthcheck is not running - nothing to stop",
	)

	err = wrapHealth.Start()
	require.NoError(t, err)
	err = wrapHealth.Stop()
	require.NoError(t, err)
}

func TestHealthChecker_AddCheck(t *testing.T) {
	wrapHealth := health.New()

	err := wrapHealth.AddURL("foo", "http://localhost/somerandomurl", 2*time.Second)
	require.NoError(t, err)

	err = wrapHealth.AddURL("", "http://localhost/somerandomurl", -1*time.Second)
	assert.EqualError(
		t,
		err,
		"The name can't be empty",
	)

	err = wrapHealth.AddURL("bar", "http://localhost/somerandomurl", -1*time.Second)
	assert.EqualError(
		t,
		err,
		"The timeout value must be greather than zero",
	)

	err = wrapHealth.AddURL("foo", "", 2*time.Second)
	assert.EqualError(
		t,
		err,
		"The target url can't be empty",
	)

	err = wrapHealth.AddURL("foobar", ":foo", 2*time.Second)
	assert.EqualError(
		t,
		err,
		"parse :foo: missing protocol scheme",
	)
}

type TestHealthChecker struct {
	fooChecked bool
	barChecked bool
	msg        chan string
}

func (thc *TestHealthChecker) UrlHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path[1:] == "foo" {
		thc.fooChecked = true
		w.WriteHeader(500)
	}
	if r.URL.Path[1:] == "bar" {
		thc.barChecked = true
	}
	if thc.barChecked == true && thc.fooChecked == true {
		thc.msg <- "foo and bar checked"
	}
}

func TestHealthChecker_State(t *testing.T) {
	var thc TestHealthChecker
	thc.msg = make(chan string, 1)
	http.HandleFunc("/foo", thc.UrlHandler)
	http.HandleFunc("/bar", thc.UrlHandler)
	go func() {
		err := http.ListenAndServe("localhost:63333", nil)
		require.NoError(t, err)
	}()

	wrapHealth := health.New()

	err := wrapHealth.AddURL("bar", "http://localhost:63333/bar", 1*time.Second)
	require.NoError(t, err)

	err = wrapHealth.State("foo")
	assert.EqualError(t, err, "HealthCheck entry 'foo' not found")

	err = wrapHealth.AddURL("foo", "http://localhost:63333/foo", 1*time.Second)
	require.NoError(t, err)

	err = wrapHealth.Start()
	require.NoError(t, err)

	select {
	case <-thc.msg:
		time.Sleep(500 * time.Millisecond)
	case <-time.After(10 * time.Second):
		t.Fatal("timeout has reached")
	}

	err = wrapHealth.State("foo")
	assert.Contains(t, err, "Received status code '500'")

	err = wrapHealth.State("bar")
	require.NoError(t, err)
}
