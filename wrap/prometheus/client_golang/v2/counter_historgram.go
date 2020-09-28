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

package v2

import (
	of "github.com/cisco-cx/of/pkg/v2"
	promclient "github.com/prometheus/client_golang/prometheus"
)

// HistogramVec represents the options required for prometheus.HistogramVec
// and reference to the created prometheus.HistogramVec.
type HistogramVec struct {
	Namespace string
	Name      string
	Help      string
	histVec   *promclient.HistogramVec
}

// Create a new histogram vec.
func (c *HistogramVec) Create(labels []string) error {
	c.histVec = promclient.NewHistogramVec(promclient.HistogramOpts{
		Namespace: c.Namespace,
		Name:      c.Name,
		Help:      c.Help,
		Buckets:   promclient.DefBuckets,
	}, labels)

	if c.histVec == nil {
		return of.ErrHistogramVecCreateFailed
	}
	return promclient.Register(c.histVec)
}

// Increment histogram for given label and value.
func (c *HistogramVec) Observed(values []string, timeTaken float64) error {
	c.histVec.WithLabelValues(values...).Observe(timeTaken)
	return nil
}
