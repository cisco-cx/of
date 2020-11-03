package cmd_test

import (
	"io/ioutil"
	"testing"
	"time"

	of "github.com/cisco-cx/of/pkg/v1"
	aci "github.com/cisco-cx/of/wrap/aci/v1"
	acigo "github.com/cisco-cx/of/wrap/acigo/v1"
	alertmanager "github.com/cisco-cx/of/wrap/alertmanager/v1"
	http "github.com/cisco-cx/of/wrap/http/v1"
	logger "github.com/cisco-cx/of/wrap/logrus/v1"
	"github.com/stretchr/testify/require"
	"gotest.tools/assert"
)

// Test handler.run
func TestHandlerRun(t *testing.T) {
	cfg := &of.ACIConfig{}
	cfg.Application = "testing_aci"
	cfg.ListenAddress = "127.0.0.1:9011"
	cfg.CycleInterval = 10
	cfg.AmURL = "locahost:9093"
	cfg.ACIHosts = []string{"::1"}
	cfg.User = "user"
	cfg.Pass = "pass"

	cfg.AlertsCFGFile = "test/alerts.yaml"
	cfg.SecretsCFGFile = "test/secrets.yaml"

	var err error
	log := logger.New()
	client, err := acigo.NewACIClient(of.ACIClientConfig{User: cfg.User, Pass: cfg.Pass}, log)
	require.NoError(t, err)
	handler := *&aci.Handler{Config: cfg, Log: log, Aci: client}
	handler.Ams = &alertmanager.AlertService{AmURL: cfg.AmURL, Version: cfg.Version}
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
