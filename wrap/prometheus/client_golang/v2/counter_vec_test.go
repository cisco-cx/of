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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	of "github.com/cisco-cx/of/pkg/v2"
	promclient "github.com/cisco-cx/of/wrap/prometheus/client_golang/v2"
)

// Ensures Counter implements of.Counter
func TestCounterVecInterface(t *testing.T) {
	var _ of.CounterVector = &promclient.CounterVec{}
}

// Ensure Counter can be incremented.
func TestCounterVecIncr(t *testing.T) {
	cntrVec := promclient.CounterVec{Namespace: "TestAppVec", Name: "test_counter_incr", Help: "This is a test counter."}
	err := cntrVec.Create([]string{"request"})
	require.NoError(t, err)
	for i := 0; i < 10; i++ {
		err = cntrVec.Incr("request", "get")
		require.NoError(t, err)
	}

	for i := 0; i < 20; i++ {
		err = cntrVec.Incr("request", "post")
		require.NoError(t, err)
	}

	// Search metrics to check if value of counter is 10.
	require.Contains(t, promMetrics(t), "TestAppVec_test_counter_incr{request=\"get\"} 10")
	require.Contains(t, promMetrics(t), "TestAppVec_test_counter_incr{request=\"post\"} 20")
	require.Panics(t, assert.PanicTestFunc(func() { cntrVec.Incr("", "") }))
}
