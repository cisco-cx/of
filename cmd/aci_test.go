package cmd_test

import (
	"io/ioutil"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"gotest.tools/assert"
	"github.com/cisco-cx/of/cmd"
	of "github.com/cisco-cx/of/lib/v1alpha1"
	aci "github.com/cisco-cx/of/wrap/aci/v1alpha1"
	acigo "github.com/cisco-cx/of/wrap/acigo/v1alpha1"
	alertmanager "github.com/cisco-cx/of/wrap/alertmanager/v1alpha1"
	http "github.com/cisco-cx/of/wrap/http/v1alpha1"
)

type DNSEntry struct {
	Hostname string
	Address  string
	Result   bool
}

// Test handler.run
func TestHandlerRun(t *testing.T) {
	cfg := &of.ACIConfig{}
	cfg.Application = "testing_aci"
	cfg.ListenAddress = "127.0.0.1:9100"
	cfg.CycleInterval = 10
	cfg.AmURL = "locahost:9093"
	cfg.ACIHost = "::1"

	cfg.AlertsCFGFile = "test/alerts.yaml"
	cfg.SecretsCFGFile = "test/secrets.yaml"

	handler := *&aci.Handler{Config: cfg}
	handler.Aci = &acigo.ACIService{cfg}
	handler.Ams = &alertmanager.AlertService{cfg}
	go handler.Run()

	time.Sleep(time.Second)
	c := http.NewClient()
	res, err := c.Get("http://" + cfg.ListenAddress + "/metrics")
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	all, err := ioutil.ReadAll(res.Body)
	require.NoError(t, err)
	require.Contains(t, string(all), "testing_aci_notification_cycle_count 1")
	handler.Shutdown()
}

// Test DNS lookup.
func TestVerifiedHost(t *testing.T) {
	entries := []DNSEntry{
		{Hostname: "google.com", Address: "fe80::800:27ff:fe00:1", Result: false},
		{Hostname: "www1.cisco.com.", Address: "2001:420:1101:1::a", Result: true},
		{Hostname: "edge-star-mini6-shv-01-sjc3.facebook.com.", Address: "2a03:2880:f131:83:face:b00c:0:25de", Result: true},
		{Hostname: "localhost", Address: "::1", Result: true},
	}
	for _, entry := range entries {
		hostname, ip := cmd.VerifiedHost(entry.Address)
		if (ip == entry.Address && hostname == entry.Hostname) != entry.Result {
			require.EqualValues(t, entry.Hostname, hostname)
			require.EqualValues(t, entry.Address, ip)
		}
	}
}