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

package v2

import (
	"fmt"
	"net/url"
	"time"

	"github.com/InVisionApp/go-health"
	"github.com/InVisionApp/go-health/checkers"
	of "github.com/cisco-cx/of/pkg/v1"
)

type HealthChecker struct {
	ih health.IHealth
}

func New() *HealthChecker {
	ih := health.New()
	return &HealthChecker{
		ih: ih,
	}
}

// Start the health checker.
func (hc *HealthChecker) Start() error {
	return hc.ih.Start()
}

func (hc *HealthChecker) Stop() error {
	return hc.ih.Stop()
}

func (hc *HealthChecker) AddURL(name string, urlTarget string, timeout time.Duration) error {
	if name == "" {
		return of.Error("The name can't be empty")
	}
	if urlTarget == "" {
		return of.Error("The target url can't be empty")
	}
	if timeout <= 0 {
		return of.Error("The timeout value must be greather than zero")
	}

	urlTargetParsed, err := url.Parse(urlTarget)
	if err != nil {
		return err
	}

	httpChecker, err := checkers.NewHTTP(&checkers.HTTPConfig{
		URL: urlTargetParsed,
	})
	if err != nil {
		return err
	}

	healthConfig := health.Config{
		Name:     name,
		Fatal:    true,
		Interval: timeout * time.Second,
		Checker:  httpChecker,
	}
	return hc.ih.AddCheck(&healthConfig)
}

func (hc *HealthChecker) State(name string) error {
	states, _, err := hc.ih.State()
	if err != nil {
		return err
	}

	state, ok := states[name]
	if !ok {
		return of.Error(fmt.Sprintf("HealthCheck entry '%s' not found", name))
	}
	if state.Status == "failed" {
		return of.Error(state.Err)
	}
	return nil
}
