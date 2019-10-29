package v1_test

import (
	"fmt"
	"io/ioutil"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"gotest.tools/assert"
	of "github.com/cisco-cx/of/pkg/v1"
	http "github.com/cisco-cx/of/wrap/http/v1"
)

// Test http.Get
func TestGet(t *testing.T) {
	server_addr := "localhost:54941"
	srv := startServer(t, server_addr)
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
func startServer(t *testing.T, server_addr string) *http.Server {
	response_text := "HandleFunc called."

	c := &of.HTTPConfig{ListenAddress: server_addr}

	srv := http.NewServer(c)
	srv.HandleFunc("/", func(w of.ResponseWriter, r of.Request) {
		fmt.Fprint(w, response_text)
	})
	go func() {
		err := srv.ListenAndServe()
		require.NoError(t, err)
	}()
	time.Sleep(time.Second)
	return srv
}

// Test client.Do request.
func TestDo(t *testing.T) {
	server_addr := "localhost:54942"
	srv := startServer(t, server_addr)
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
