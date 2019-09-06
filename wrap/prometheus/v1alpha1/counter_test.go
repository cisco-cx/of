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

package v1alpha1_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/stretchr/testify/require"
	prometheus "github.com/cisco-cx/of/wrap/prometheus/v1alpha1"
)

// Ensure Counter can be created.
func TestCounter_Create(t *testing.T) {
	cntr := prometheus.Counter{Namespace: "TestApp", Name: "test_counter_create", Help: "This is a test counter."}
	err := cntr.Create()
	require.NoError(t, err)
	defer cntr.Destroy()

	// Search metrics for newly created counter.
	require.Contains(t, promMetrics(t), "TestApp_test_counter_create 0")
}

// Ensure Counter can be incremented.
func TestCounter_Incr(t *testing.T) {
	cntr := prometheus.Counter{Namespace: "TestApp", Name: "test_counter_incr", Help: "This is a test counter."}
	err := cntr.Create()
	require.NoError(t, err)
	defer cntr.Destroy()
	for i := 0; i < 10; i++ {
		err = cntr.Incr()
		require.NoError(t, err)
	}

	// Search metrics to check if value of counter is 10.
	require.Contains(t, promMetrics(t), "TestApp_test_counter_incr 10")
}

// Ensure Counter can be destroyed.
func TestCounter_Destroy(t *testing.T) {
	cntr := prometheus.Counter{Namespace: "TestApp", Name: "test_counter_destroy", Help: "This is a test counter."}
	err := cntr.Create()
	require.NoError(t, err)
	err = cntr.Destroy()
	require.NoError(t, err)

	// Search metrics to check if counter has been removed.
	require.NotContains(t, promMetrics(t), "TestApp_test_counter_destroy")
}

// Fetches current metrics.
func promMetrics(t *testing.T) string {
	ts := httptest.NewServer(promhttp.Handler())
	defer ts.Close()
	res, err := http.Get(ts.URL)
	require.NoError(t, err)
	defer res.Body.Close()
	metrics, err := ioutil.ReadAll(res.Body)
	require.NoError(t, err)
	return string(metrics)
}
