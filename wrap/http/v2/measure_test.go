package v2_test

import (
	"io/ioutil"
	go_http "net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	of "github.com/cisco-cx/of/pkg/v2"
	http "github.com/cisco-cx/of/wrap/http/v2"
	promclient "github.com/cisco-cx/of/wrap/prometheus/client_golang/v2"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/stretchr/testify/require"
)

// Enforce interface implementation.
func TestMeasurerInterface(t *testing.T) {
	var _ of.Measurer = &http.Measurer{}
}

// Test of.Measurer's Measure method
func TestMeasure(t *testing.T) {

	histVec := promclient.HistogramVec{Namespace: t.Name(), Name: "Test", Help: "Testing measure."}
	err := histVec.Create([]string{"method", "uri", "status_code"})
	require.NoError(t, err)

	mr := http.Measurer{&histVec}

	fn := func(rw of.ResponseWriter, r of.Request) {
		time.Sleep(500 * time.Millisecond)
		rw.WriteHeader(403)
	}

	w := httptest.NewRecorder()
	rw := http.NewResponseWriter(w)
	var r of.Request = &go_http.Request{
		Method: "GET",
		URL: &url.URL{
			Host: "localhost",
			Path: "/foobar",
		},
	}

	mr.Measure(rw, r, fn)
	metrics := promMetrics(t)
	require.Contains(t, metrics, `TestMeasure_Test_bucket{method="GET",status_code="403",uri="/foobar",le="1"} 1`)
}

// Fetches current metrics.
func promMetrics(t *testing.T) string {
	ts := httptest.NewServer(promhttp.Handler())
	defer ts.Close()
	res, err := go_http.Get(ts.URL)
	require.NoError(t, err)
	defer res.Body.Close()
	metrics, err := ioutil.ReadAll(res.Body)
	require.NoError(t, err)
	return string(metrics)
}
