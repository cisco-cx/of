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
	promclient "github.com/prometheus/client_golang/prometheus"
	of "github.com/cisco-cx/of/pkg/v2"
)

// CounterVec represents the options required for prometheus.CounterVec
// and reference to the created prometheus.CounterVec.
type CounterVec struct {
	Namespace string
	Name      string
	Help      string
	cntrVec   *promclient.CounterVec
}

// Create a new counter vec.
func (c *CounterVec) Create(labels []string) error {
	c.cntrVec = promclient.NewCounterVec(promclient.CounterOpts{
		Namespace: c.Namespace,
		Name:      c.Name,
		Help:      c.Help,
	}, labels)

	if c.cntrVec == nil {
		return of.ErrCounterVecCreateFailed
	}
	return promclient.Register(c.cntrVec)
}

// Increment counter for given label and value.
func (c *CounterVec) Incr(label string, value string) error {
	c.cntrVec.With(promclient.Labels{label: value}).Inc()
	return nil
}
