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
	"path/filepath"
	"strconv"
	"time"

	of "github.com/cisco-cx/of/pkg/v2"
	promclient "github.com/cisco-cx/of/wrap/prometheus/client_golang/v2"
)

// Implements of.Measurer
type Measurer struct {
	HistVec *promclient.HistogramVec
}

// Method to measure http response as a prometheus metric.
func (m *Measurer) Measure(rw of.ResponseWriter, r of.Request, fn func(rw of.ResponseWriter, r of.Request)) {
	t1 := time.Now()
	fn(rw, r)
	t2 := time.Now().Sub(t1).Seconds()

	values := []string{
		r.Method,
		filepath.Join(r.Host, r.URL.Path),
		strconv.Itoa(rw.StatusCode()),
	}
	m.HistVec.Observed(values, t2)
}
