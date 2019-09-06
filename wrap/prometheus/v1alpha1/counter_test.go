package v1alpha1

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/stretchr/testify/require"
)

// Ensure Counter can be created.
func TestCounter_Create(t *testing.T) {
	cntr := Counter{Namespace: "TestApp", Name: "test_counter_create", Help: "This is a test counter."}
	err := cntr.Create()
	require.NoError(t, err)
	defer cntr.Destroy()

	// Search metrics for newly created counter.
	require.Contains(t, promMetrics(t), "TestApp_test_counter_create 0")
}

// Ensure Counter can be incremented.
func TestCounter_Incr(t *testing.T) {
	cntr := Counter{Namespace: "TestApp", Name: "test_counter_incr", Help: "This is a test counter."}
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

func TestCounter_Destroy(t *testing.T) {
	cntr := Counter{Namespace: "TestApp", Name: "test_counter_destroy", Help: "This is a test counter."}
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
