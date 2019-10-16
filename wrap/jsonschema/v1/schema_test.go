package v1_test

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/require"
	of "github.com/cisco-cx/of/pkg/v1"
	js "github.com/cisco-cx/of/wrap/jsonschema/v1"
)

var schemaBytes = []byte(`{
  "$id": "https://example.com/person.schema.json",
  "$schema": "http://json-schema.org/draft-07/schema#",
  "title": "Person",
  "type": "object",
  "properties": {
    "firstName": {
      "type": "string",
      "description": "The person's first name."
    },
    "lastName": {
      "type": "string",
      "description": "The person's last name."
    },
    "age": {
      "description": "Age in years which must be equal to or greater than zero.",
      "type": "integer",
      "minimum": 0
    }
  }
}`)

// Enforce implementation of of/Schema
func TestInterface(t *testing.T) {
	var _ of.SchemaLoader = &js.Schema{}
	var _ of.SchemaValidator = &js.Schema{}
}

// Test schema loader
func TestLoader(t *testing.T) {
	_ = newSchema(t)
}

// Test schema validator with a JSON input
func TestValidateJSON(t *testing.T) {
	schema := newSchema(t)
	jsonBytes := []byte(`{
							  "firstName": "John",
							  "lastName": "Doe",
							  "age": 21
						}`)

	err := schema.ValidateJSON(jsonBytes)
	require.NoError(t, err)
}

// Test schema validator with a JSON input
func TestValidateYAML(t *testing.T) {
	schema := newSchema(t)
	yamlBytes := []byte(`
firstName: John
lastName: Doe
age: 21`)

	err := schema.ValidateYAML(yamlBytes)
	require.NoError(t, err)
}

// Returns a *js.Schema with test schema loaded.
func newSchema(t *testing.T) *js.Schema {
	schema := &js.Schema{}
	err := schema.Load(schemaBytes)
	require.NoError(t, err)
	return schema
}

// Test ACI Alerts schema.
func TestACIAlertsSchema(t *testing.T) {
	testSchema(t, "schema/aci/alerts.schema", "sample/aci/alerts.yaml")
}

// Test ACI Secrets schema.
func TestACISecretsSchema(t *testing.T) {
	testSchema(t, "schema/aci/secrets.schema", "sample/aci/secrets.yaml")
}

// Test SNMP Alerts schema.
func TestSNMPAlertsSchema(t *testing.T) {
	testSchema(t, "schema/snmp/alerts.schema", "sample/snmp/alerts.yaml")
}

// Test Alerts schema.
func testSchema(t *testing.T, schemaFile string, yamlFile string) {
	schemaContent, err := ioutil.ReadFile(schemaFile)
	require.NoError(t, err)
	yaml, err := ioutil.ReadFile(yamlFile)
	require.NoError(t, err)

	schema := &js.Schema{}
	err = schema.Load(schemaContent)
	require.NoError(t, err)

	err = schema.ValidateYAML(yaml)
	require.NoError(t, err)
}
