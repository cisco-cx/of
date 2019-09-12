package v1alpha1_test

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/require"
	"gotest.tools/assert"
	of "github.com/cisco-cx/of/lib/v1alpha1"
	http "github.com/cisco-cx/of/wrap/http/v1alpha1"
)

var server_addr string = "localhost:54931"

// Test http.Get
func TestGet(t *testing.T) {
	srv := startServer(t)
	defer srv.Shutdown()
	c := http.NewClient()
	res, err := c.Get("http://" + server_addr)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	all, err := ioutil.ReadAll(res.Body)
	require.NoError(t, err)
	assert.Equal(t, "HandleFunc called.", string(all))
}

// Start a HTTP Server to test.
func startServer(t *testing.T) *http.Server {
	response_text := "HandleFunc called."
	server := of.Server{
		Addr: server_addr,
	}
	srv := http.NewServer(server)
	srv.HandleFunc("/", func(w of.ResponseWriter, r of.Request) {
		fmt.Fprint(w, response_text)
	})
	go func() {
		err := srv.ListenAndServe()
		require.NoError(t, err)
	}()
	return srv
}

// Test client.Do request.
func TestDo(t *testing.T) {
	srv := startServer(t)
	defer srv.Shutdown()
	c := http.NewClient()

	req, err := http.NewRequest("Get", "http://"+server_addr, nil)
	require.NoError(t, err)
	req.Header.Set("User-Agent", "test")
	req.Header.Set("Content-Type", "application/json")

	res, err := c.Do(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	all, err := ioutil.ReadAll(res.Body)
	require.NoError(t, err)
	assert.Equal(t, "HandleFunc called.", string(all))
}