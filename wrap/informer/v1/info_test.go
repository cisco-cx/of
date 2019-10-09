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
package v1_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/stretchr/testify/require"

	of "github.com/cisco-cx/of/pkg/v1"
	informer "github.com/cisco-cx/of/wrap/informer/v1"
)

// Confirm that informer.InfoService implements the of.InfoService interface.
func TestInfoService_Interface(t *testing.T) {
	var _ of.InfoService = &informer.InfoService{}
	require.Nil(t, nil) // If we get this far, the test passed.
}

// Confirm that informer.InfoService implements the of.InfoCollector interface.
func TestInfoCollector_Interface(t *testing.T) {
	var _ of.InfoCollector = &informer.InfoService{}
	require.Nil(t, nil) // If we get this far, the test passed.
}

// Confirm BuildInfo() of.InfoService functionality of informer.InfoService.
func TestInfoService_BuildInfo(t *testing.T) {
	s := getMockInfoService()
	expect := "(language=language, languageVersion=languageVersion, user=user, date=date)"
	require.Equal(t, expect, s.BuildInfo())
}

// Confirm Metadata() of.InfoService functionality of informer.InfoService.
func TestInfoService_Metadata(t *testing.T) {
	s := getMockInfoService()
	expect := "(program=program, license=license, url=url)"
	require.Equal(t, expect, s.Metadata())
}

// Confirm String() of.InfoService functionality of informer.InfoService.
func TestInfoService_String(t *testing.T) {
	s := getMockInfoService()
	expect := "(metadata=(program=program, license=license, url=url), versionInfo=(version=version, branch=branch, revision=revision), buildInfo=(language=language, languageVersion=languageVersion, user=user, date=date))"
	require.Equal(t, expect, s.String())
}

// Confirm VersionInfo() of.InfoService functionality of informer.InfoService.
func TestInfoService_VersionInfo(t *testing.T) {
	s := getMockInfoService()
	expect := "(version=version, branch=branch, revision=revision)"
	require.Equal(t, expect, s.VersionInfo())
}

// Confirm basic of.InfoCollector functionality of informer.InfoService.
func TestInfoCollector_Basic(t *testing.T) {
	// Prepare to require multiple times.
	require := require.New(t)

	// Get a new service.
	s := getMockInfoService()

	// Register the InfoService collector with the Prometheus client.
	err := s.Register()
	require.NoError(err)
	defer s.Unregister()

	// Parse metrics.
	require.Contains(getMetrics(t), `program_info{branch="branch",build_date="date",build_user="user",language="language",language_version="languageVersion",license="license",program="program",revision="revision",url="url",version="version"} 1`)
}

// getMetrics bootstraps an http server and fetches current metrics.
func getMetrics(t *testing.T) string {
	ts := httptest.NewServer(promhttp.Handler())
	defer ts.Close()
	res, err := http.Get(ts.URL)
	require.NoError(t, err)
	defer res.Body.Close()
	metrics, err := ioutil.ReadAll(res.Body)
	require.NoError(t, err)
	return string(metrics)
}

// getMockInfoService returns a mock *informer.InfoService.
func getMockInfoService() *informer.InfoService {
	return informer.NewInfoService(
		"program",
		"license",
		"url",
		"user",
		"date",
		"language",
		"languageVersion",
		"version",
		"revision",
		"branch",
	)
}
