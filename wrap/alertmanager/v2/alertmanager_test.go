package v2_test

import (
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	of "github.com/cisco-cx/of/pkg/v2"
	am "github.com/cisco-cx/of/wrap/alertmanager/v2"
	http "github.com/cisco-cx/of/wrap/http/v2"
	logger "github.com/cisco-cx/of/wrap/logrus/v2"
)

func TestDryRun(t *testing.T) {
	// Start fake AM Server.

	addr := "localhost:17931"
	timeout := 5 * time.Second

	hc := &of.HTTPConfig{ListenAddress: addr, ReadTimeout: timeout, WriteTimeout: timeout}

	srv := http.NewServer(hc)
	srv.HandleFunc("/api/v1/alerts", func(w of.ResponseWriter, r of.Request) {
		w.WriteHeader(http.StatusBadRequest)
	})

	err := srv.ListenAndServe()
	require.NoError(t, err)
	defer srv.Shutdown()

	// setup logger
	output := &bytes.Buffer{}
	log := logger.New()
	log.SetOutput(output)

	// Setup alert service.
	as := am.AlertService{
		AmURL:  "http://" + addr,
		Log:    log,
		DryRun: true,
	}

	alerts := &[]of.Alert{
		of.Alert{
			Labels: map[string]string{"dryRun": "true"},
		},
	}
	err = as.Notify(alerts)
	require.NoError(t, err)
	dryRunContents := string(output.Bytes())
	require.Contains(t, dryRunContents, "level=info msg=\"Dry run.\" annotations=\"map[]\" endsAt=\"0001-01-01 00:00:00 +0000 UTC\" generatorURL= labels=\"map[dryRun:true]\" startsAt=\"0001-01-01 00:00:00 +0000 UTC\"")

	// Disable dry run.
	as.DryRun = false
	err = as.Notify(alerts)
	require.Error(t, err)

}
