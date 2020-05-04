package cmd_test

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
	snmp_cmd "github.com/cisco-cx/of/cmd"
	of_v2 "github.com/cisco-cx/of/pkg/v2"
	http "github.com/cisco-cx/of/wrap/http/v2"
	logger "github.com/cisco-cx/of/wrap/logrus/v2"
	snmp "github.com/cisco-cx/of/wrap/snmp/v2"
)

// Test SNMP mibs pre process
func TestSNMPMIBsPreprocess(t *testing.T) {

	// Create temp. file
	tmpFile, err := ioutil.TempFile(os.TempDir(), "cache-")
	require.NoError(t, err)
	fileName := tmpFile.Name()
	tmpFile.Close()

	defer os.Remove(fileName)

	// Using mibs dir.
	cmd := &cobra.Command{}
	args := []string{
		"--mibs-dir=test/snmp/mibs/",
		"--cache-file=" + fileName,
	}

	snmp_cmd.RunMibsPreProcess(cmd, args)

	// Compute check sum of generated file.
	f, err := os.Open(fileName)
	require.NoError(t, err)
	defer f.Close()

	h := md5.New()
	_, err = io.Copy(h, f)
	require.NoError(t, err)
	computedHash := fmt.Sprintf("%x", h.Sum(nil))
	require.Equal(t, "c81bbe474a939174697cc9eb784df704", computedHash)

}

// Test handler.run
func TestSNMPHandler(t *testing.T) {

	amAddress := "localhost:15932"

	// Start fake AM Server.
	srv := startFakeAM(t, amAddress)
	defer srv.Shutdown()

	// Using mibs dir.
	args := []string{
		"--listen-address=blah",
		"--listen-address=localhost:14444",
		"--am-address=http://" + amAddress,
		"--mibs-dir=test/snmp/mibs/",
		"--config-dir=test/snmp/configs/",
	}
	checkHandler(t, args, amAddress, "test_dir")

	// Using cache file that was just generated dir.
	args = []string{
		"--listen-address=blah",
		"--listen-address=localhost:14444",
		"--am-address=http://" + amAddress,
		"--am-timeout=10s",
		"--cache-file=test/snmp/cache_mibs.json",
		"--config-dir=test/snmp/configs/",
	}
	checkHandler(t, args, amAddress, "test_cache")
}

func checkHandler(t *testing.T, args []string, amAddress string, cn string) {
	cmd := &cobra.Command{}
	snmp_cmd.ParseSNMPHandlerFlags(cmd, args)
	config := snmp_cmd.SNMPConfig(cmd)
	config.Application = cn

	log := logger.New()
	cntr, cntrVec := snmp.InitCounters(config.Application, log)
	service, err := snmp.NewService(log, config, cntr, cntrVec)
	require.NoError(t, err)

	handler := &snmp.Handler{
		Config: config,
		SNMP:   service,
		Log:    log,
	}

	handler.Run()

	runHandlerChecks(t, config.ListenAddress)

	handler.Shutdown()
}

func startFakeAM(t *testing.T, amAddress string) *http.Server {

	hc := &of_v2.HTTPConfig{ListenAddress: amAddress, ReadTimeout: 1 * time.Second, WriteTimeout: 1 * time.Second}

	srv := http.NewServer(hc)
	srv.HandleFunc("/-/healthy", func(w of_v2.ResponseWriter, r of_v2.Request) {
		w.WriteHeader(http.StatusBadGateway)
	})
	srv.HandleFunc("/api/v1/alerts", func(w of_v2.ResponseWriter, r of_v2.Request) {

		var alerts []of_v2.Alert
		err := json.NewDecoder(r.Body).Decode(&alerts)
		require.NoError(t, err)
		expectedOIDStr := "iso.org.dod.internet.private.enterprises.starentMIB.starentTraps.starAAAAccServerMisconfigured"
		for _, alert := range alerts {
			for k, v := range alert.Annotations {
				if k == "event_oid_string" {
					require.Equal(t, expectedOIDStr, v)
				}
			}
		}
	})

	err := srv.ListenAndServe()
	require.NoError(t, err)
	return srv
}

func runHandlerChecks(t *testing.T, listenAddress string) {

	// Test metrics
	message := getResponse(t, 200, "http://"+listenAddress+"/metrics")
	require.Contains(t, message, "promhttp_metric_handler_requests_total")

	// Test Status
	checkResponse(t, 200, "http://"+listenAddress+"/api/v2/status", "{ AlertManager Client for SNMP Traps {https://github.com/cisco-cx/am-client-snmp} success}")

	// Test Posting SNMP event
	dataBytes, err := json.Marshal(TrapEvents())
	require.NoError(t, err)

	data := bytes.NewBuffer(dataBytes)
	c := http.NewClient()
	req, err := http.NewRequest("Post", "http://"+listenAddress+"/api/v2/events", data)
	require.NoError(t, err)
	req.Header.Set("User-Agent", "test")
	req.Header.Set("Content-Type", "application/json")

	res, err := c.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)
}

// HTTP client to hit server and check response.
func checkResponse(t *testing.T, expectedStatusCode int, u, msg string) {
	all := getResponse(t, expectedStatusCode, u)
	require.Equal(t, msg, all)
}

// HTTP client to hit server and check response.
func getResponse(t *testing.T, expectedStatusCode int, u string) string {

	c := http.NewClient()
	res, err := c.Get(u)
	require.NoError(t, err)
	require.Equal(t, expectedStatusCode, res.StatusCode)

	all, err := ioutil.ReadAll(res.Body)
	require.NoError(t, err)
	return string(all)
}
