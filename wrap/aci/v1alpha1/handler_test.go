package v1alpha1_test

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/require"
	of "github.com/cisco-cx/of/lib/v1alpha1"
	aci "github.com/cisco-cx/of/wrap/aci/v1alpha1"
	mapstructure "github.com/cisco-cx/of/wrap/mapstructure/v1alpha1"
)

const faultJson = "test/faultInst.json"

// Test FaultToAlerts.
func TestFaultToAlerts(t *testing.T) {
	c := &of.ACIConfig{}
	c.Application = "testing_handler"
	c.AlertsCFGFile = "test/alerts.yaml"
	c.SecretsCFGFile = "test/secrets.yaml"
	handler := &aci.Handler{Config: c}
	handler.InitHandler()
	faults, err := handler.FaultsToAlerts(getFaults(t))
	require.NoError(t, err)
	faults_json, err := json.Marshal(faults)
	require.NoError(t, err)

	hash := md5.New()
	_, err = hash.Write(faults_json)
	require.NoError(t, err)
	md5sum := fmt.Sprintf("%x", hash.Sum(nil))
	require.Equal(t, "75d506f347276107b1fe7657b681473b", md5sum)
}

// Returns faults in faultJson file.
func getFaults(t *testing.T) []of.Map {
	faults, err := ioutil.ReadFile(faultJson)
	require.NoError(t, err)
	list, err := jsonImdataAttributes(faults, "faultInst", "FaultList")
	require.NoError(t, err)

	mm := make([]of.Map, len(list))
	for i, v := range list {
		//fmt.Printf("i : %v, v : %v", i, v)
		mapstructure.NewMap(v).DecodeMap(&mm[i])
	}
	fmt.Sprintf("%+v\n", mm)
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
