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

package v2_test

import (
	"testing"

	of "github.com/cisco-cx/of/pkg/v2"
	promclient "github.com/cisco-cx/of/wrap/prometheus/client_golang/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Ensures HistogramVec implements of.HistogramVector
func TestHistogramVecInterface(t *testing.T) {
	var _ of.HistogramVector = &promclient.HistogramVec{}
}

// Test Histogram.
func TestHistogram(t *testing.T) {
	histVec := promclient.HistogramVec{Namespace: "TestAppHistVec", Name: "test_histogram_vec", Help: "This is a test histogram."}
	err := histVec.Create([]string{"method", "uri", "status_code"})
	require.NoError(t, err)
	histVec.Observed([]string{"GET", "/foo", "200"}, 10)
	histVec.Observed([]string{"GET", "/foo", "200"}, 1)
	histVec.Observed([]string{"GET", "/foo", "200"}, 4)
	histVec.Observed([]string{"GET", "/bar", "500"}, 8)
	histVec.Observed([]string{"GET", "/foo", "403"}, 10)

	// Search metrics to check for histogram values..
	require.Contains(t, promMetrics(t), `TestAppHistVec_test_histogram_vec_bucket{method="GET",status_code="200",uri="/foo",le="1"} 1`)
	require.Contains(t, promMetrics(t), `TestAppHistVec_test_histogram_vec_bucket{method="GET",status_code="200",uri="/foo",le="2.5"} 1`)
	require.Contains(t, promMetrics(t), `TestAppHistVec_test_histogram_vec_bucket{method="GET",status_code="200",uri="/foo",le="5"} 2`)
	require.Contains(t, promMetrics(t), `TestAppHistVec_test_histogram_vec_bucket{method="GET",status_code="200",uri="/foo",le="10"} 3`)
	require.Contains(t, promMetrics(t), `TestAppHistVec_test_histogram_vec_bucket{method="GET",status_code="200",uri="/foo",le="+Inf"} 3`)

	require.Contains(t, promMetrics(t), `TestAppHistVec_test_histogram_vec_bucket{method="GET",status_code="403",uri="/foo",le="10"} 1`)
	require.Contains(t, promMetrics(t), `TestAppHistVec_test_histogram_vec_bucket{method="GET",status_code="403",uri="/foo",le="+Inf"} 1`)

	require.Contains(t, promMetrics(t), `TestAppHistVec_test_histogram_vec_bucket{method="GET",status_code="500",uri="/bar",le="10"} 1`)
	require.Contains(t, promMetrics(t), `TestAppHistVec_test_histogram_vec_bucket{method="GET",status_code="500",uri="/bar",le="+Inf"} 1`)
	require.Panics(t, assert.PanicTestFunc(func() { histVec.Observed([]string{"", ""}, 0) }))
}
