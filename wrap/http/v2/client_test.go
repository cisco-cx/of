package v2_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"testing"

	of "github.com/cisco-cx/of/pkg/v2"
	http "github.com/cisco-cx/of/wrap/http/v2"
	"github.com/stretchr/testify/require"
	"gotest.tools/assert"
)

// Test http.Get
func TestGet(t *testing.T) {
	server_addr := "localhost:64941"
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

	srv := http.NewServer(c, t.Name())
	srv.HandleFunc("/", func(w of.ResponseWriter, r of.Request) {
		fmt.Fprint(w, response_text)
	})
	srv.HandleFunc("/posttest", func(w of.ResponseWriter, r of.Request) {
		data := make(map[string]string)
		err := json.NewDecoder(r.Body).Decode(&data)
		require.NoError(t, err)
		expectedData := map[string]string{
			"data": "This is a post request.",
		}
		require.Equal(t, expectedData, data)
		fmt.Fprint(w, "Post test called.")
	})
	err := srv.ListenAndServe()
	require.NoError(t, err)
	return srv
}

// Test client.Do request.
func TestDo(t *testing.T) {
	server_addr := "localhost:64942"
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

// Test client.Do request.
func TestDoPost(t *testing.T) {
	server_addr := "localhost:64942"
	srv := startServer(t, server_addr)
	defer srv.Shutdown()
	c := http.NewClient()
	data := strings.NewReader(`{"data":"This is a post request."}`)
	req, err := http.NewRequest("Post", "http://"+server_addr+"/posttest", data)
	require.NoError(t, err)
	req.Header.Set("User-Agent", "test")
	req.Header.Set("Content-Type", "application/json")

	res, err := c.Do(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	all, err := ioutil.ReadAll(res.Body)
	require.NoError(t, err)
	assert.Equal(t, "Post test called.", string(all))
}
