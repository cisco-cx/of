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

package v1alpha1

import (
	promclient "github.com/prometheus/client_golang/prometheus"
	of "github.com/cisco-cx/of/lib/v1alpha1"
)

// Ensures Counter implements of.Counter
var _ of.Counter = &Counter{}

// Counter represents the options required for promclient.Counter
// and reference to the created promclient.Counter.
type Counter struct {
	Namespace string
	Name      string
	Help      string
	cntr      promclient.Counter
}

// Create a new counter.
func (c *Counter) Create() error {
	c.cntr = promclient.NewCounter(promclient.CounterOpts{
		Namespace: c.Namespace,
		Name:      c.Name,
		Help:      c.Help,
	})

	if c.cntr == nil {
		return of.ErrCounterCreateFailed
	}

	return promclient.Register(c.cntr)
}

// Increment the counter by 1.
func (c *Counter) Incr() error {
	c.cntr.Inc()
	return nil
}

// Remove counter..
func (c *Counter) Destroy() error {
	if promclient.Unregister(c.cntr) {
		return nil
	}
	return of.ErrCounterDestroyFailed
}
