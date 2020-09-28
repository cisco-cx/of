package v2_test

import (
	"fmt"
	"io/ioutil"
	"testing"

	of "github.com/cisco-cx/of/pkg/v2"
	http "github.com/cisco-cx/of/wrap/http/v2"
	"github.com/stretchr/testify/require"
)

// Implementation of of.Handler to test http.Handle
type testHandler struct {
}

func (s *testHandler) ServeHTTP(rw of.ResponseWriter, r of.Request) {
	rw.Write([]byte("This is a handler."))
}

// Enforce interface implementation.
func TestInterface(t *testing.T) {
	var _ of.Serve = &http.Server{}
}

// Test server start and shutdown.
func TestServer(t *testing.T) {

	addr := "localhost:64931"

	c := &of.HTTPConfig{ListenAddress: addr}

	srv := http.NewServer(c, t.Name())

	err := srv.ListenAndServe()
	require.NoError(t, err)

	err = srv.Shutdown()
	require.NoError(t, err)
}

// Test http.HandleFunc
func TestHandleFunc(t *testing.T) {
	response_text := "HandleFunc called."
	addr := "localhost:64932"
	c := &of.HTTPConfig{ListenAddress: addr}

	srv := http.NewServer(c, t.Name())
	srv.HandleFunc("/", func(w of.ResponseWriter, r of.Request) {
		fmt.Fprint(w, response_text)
	})

	err := srv.ListenAndServe()
	require.NoError(t, err)
	checkResponse(t, "http://"+addr, response_text)
	err = srv.Shutdown()
	require.NoError(t, err)
}

// Test http.Handle
func TestHandle(t *testing.T) {
	addr := "localhost:64933"
	c := &of.HTTPConfig{ListenAddress: addr}

	srv := http.NewServer(c, t.Name())
	srv.Handle("/", &testHandler{})

	err := srv.ListenAndServe()
	require.NoError(t, err)
	checkResponse(t, "http://"+addr, "This is a handler.")
	err = srv.Shutdown()
	require.NoError(t, err)
}

// HTTP client to hit server and check response.
func checkResponse(t *testing.T, u, msg string) {

	c := http.NewClient()
	res, err := c.Get(u)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)

	all, err := ioutil.ReadAll(res.Body)
	require.NoError(t, err)
	require.Equal(t, msg, string(all))
}
