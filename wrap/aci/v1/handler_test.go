package v1_test

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"testing"

	of "github.com/cisco-cx/of/pkg/v1"
	aci "github.com/cisco-cx/of/wrap/aci/v1"
	logger "github.com/cisco-cx/of/wrap/logrus/v1"
	mapstructure "github.com/cisco-cx/of/wrap/mapstructure/v1"
	"github.com/stretchr/testify/require"
)

const faultJson = "test/faultInst.json"
const nodeJson = "test/node.json"

// Test FaultToAlerts.
func TestFaultToAlerts(t *testing.T) {
	c := &of.ACIConfig{}
	c.Application = t.Name()
	c.AlertsCFGFile = "test/alerts.yaml"
	c.SecretsCFGFile = "test/secrets.yaml"
	handler := &aci.Handler{Config: c, Log: logger.New()}
	handler.InitHandler()
	faults, err := handler.FaultsToAlerts(getFaults(t), getNodes(t))
	require.NoError(t, err)
	faults_json, err := json.Marshal(faults)
	require.NoError(t, err)

	hash := md5.New()
	_, err = hash.Write(faults_json)
	require.NoError(t, err)
	md5sum := fmt.Sprintf("%x", hash.Sum(nil))
	require.Equal(t, "2d61b4b093e4376be98d6b7aaf298ada", md5sum)
}

// Test FaultToAlerts.
func TestIdentifyNodes(t *testing.T) {
	c := &of.ACIConfig{}
	c.Application = t.Name()
	c.AlertsCFGFile = "test/alerts_identify_node.yaml"
	c.SecretsCFGFile = "test/secrets.yaml"
	handler := &aci.Handler{Config: c, Log: logger.New()}
	handler.InitHandler()
	faults, err := handler.FaultsToAlerts(getFaults(t), getNodes(t))
	require.NoError(t, err)
	faults_json, err := json.Marshal(faults)
	require.NoError(t, err)
	hash := md5.New()
	_, err = hash.Write(faults_json)
	require.NoError(t, err)
	md5sum := fmt.Sprintf("%x", hash.Sum(nil))
	// Hash of test result that handles, mapping to leaf, controller and spine,
	// as well as alerts, where topology is present, but does not map to any node.
	// Handles TLD suffix too.
	require.Equal(t, "94da3cdef371267fe04525b6889c94a9", md5sum)
}

//Test Throttle
func TestThrottle(t *testing.T) {
	c := &of.ACIConfig{}
	c.Throttle = true
	c.PostTime = 300
	c.SleepTime = 100
	c.SendTime = 30000

	count := 0
	f := func(start int, end int) {
		count += 1
		fmt.Printf("start : %d, end : %d, count : %d\n", start, end, count)
	}
	handler := &aci.Handler{Config: c, Log: logger.New()}
	handler.Throttle(15000, f)
	fmt.Printf("count : %d\n", count)
	require.Equal(t, 75, count)
}

// Returns faults in faultJson file.
func getNodes(t *testing.T) map[string]map[string]interface{} {
	data, err := ioutil.ReadFile(nodeJson)
	require.NoError(t, err)
	require.NoError(t, err)

	nodes := make([]map[string]interface{}, 0)
	err = json.Unmarshal(data, &nodes)
	require.NoError(t, err)

	nodeMap := make(map[string]map[string]interface{})
	for _, node := range nodes {
		dn := strings.Replace(node["dn"].(string), "/sys", "", -1)
		nodeMap[dn] = node
	}
	return nodeMap
}

// Returns faults in faultJson file.
func getFaults(t *testing.T) []of.Map {
	faults, err := ioutil.ReadFile(faultJson)
	require.NoError(t, err)
	list, err := jsonImdataAttributes(faults, "faultInst", "FaultList")
	require.NoError(t, err)

	mm := make([]of.Map, len(list))
	for i, v := range list {
		mapstructure.NewMap(v).DecodeMap(&mm[i])
	}
	return mm
}

// Helper function from upstream acigo, to convert faults.json into the format acigo returns.
func jsonImdataAttributes(body []byte, key, label string) ([]map[string]interface{}, error) {

	var reply interface{}
	errJSON := json.Unmarshal(body, &reply)
	if errJSON != nil {
		return nil, errJSON
	}

	return imdataAttributes(reply, key, label)
}

// Helper function from upstream acigo, to convert faults.json into the format acigo returns.
func imdataAttributes(reply interface{}, key, label string) ([]map[string]interface{}, error) {

	imdata, errImdata := mapGet(reply, "imdata")
	if errImdata != nil {
		return nil, fmt.Errorf("%s: missing imdata: %v", label, errImdata)
	}

	list, isList := imdata.([]interface{})
	if !isList {
		return nil, fmt.Errorf("%s: imdata does not hold a list", label)
	}

	return extractKeyAttributes(list, key, label), nil
}

// Helper function from upstream acigo, to convert faults.json into the format acigo returns.
func extractKeyAttributes(list []interface{}, key, label string) []map[string]interface{} {

	result := make([]map[string]interface{}, 0, len(list))

	for _, i := range list {
		item, errItem := mapGet(i, key)
		if errItem != nil {
			continue
		}
		attr, errAttr := mapGet(item, "attributes")
		if errAttr != nil {
			continue
		}
		m, isMap := attr.(map[string]interface{})
		if !isMap {
			continue
		}
		result = append(result, m)
	}

	return result
}

// Helper function from upstream acigo, to convert faults.json into the format acigo returns.
func mapGet(i interface{}, member string) (interface{}, error) {
	m, isMap := i.(map[string]interface{})
	if !isMap {
		return nil, fmt.Errorf("json mapGet: not a map")
	}
	mem, found := m[member]
	if !found {
		return nil, fmt.Errorf("json mapGet: member [%s] not found", member)
	}
	return mem, nil
}

type DNSEntry struct {
	Hostname string
	Address  string
	Result   bool
}

// Test DNS lookup.
func TestVerifiedHost(t *testing.T) {
	c := &of.ACIConfig{}
	c.Application = t.Name()
	c.AlertsCFGFile = "test/alerts.yaml"
	c.SecretsCFGFile = "test/secrets.yaml"
	handler := &aci.Handler{Config: c, Log: logger.New()}
	handler.InitHandler()
	entries := []DNSEntry{
		{Hostname: "google.com", Address: "fe80::800:27ff:fe00:1", Result: false},
		{Hostname: "www1.cisco.com.", Address: "2001:420:1101:1::a", Result: true},
		{Hostname: "edge-star-mini6-shv-01-sjc3.facebook.com.", Address: "2a03:2880:f131:83:face:b00c:0:25de", Result: true},
		{Hostname: "localhost", Address: "::1", Result: true},
	}
	for _, entry := range entries {
		hostname, ip := handler.VerifiedHost(entry.Address)
		if (ip == entry.Address && hostname == entry.Hostname) != entry.Result {
			require.EqualValues(t, entry.Hostname, hostname)
			require.EqualValues(t, entry.Address, ip)
		}
	}
}
