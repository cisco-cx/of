package v1alpha1

import (
	json_encoding "encoding/json"
	"errors"
	"fmt"

	"github.com/ghodss/yaml"
	"github.com/qri-io/jsonschema"
)

type Schema struct {
	rs *jsonschema.RootSchema
}

// Load given jsonschema.
func (j *Schema) Load(schema []byte) error {
	j.rs = &jsonschema.RootSchema{}
	return json_encoding.Unmarshal(schema, j.rs)
}

// Validate given data against loaded jsonschema.
func (j *Schema) ValidateJSON(data []byte) error {
	valErr, err := j.rs.ValidateBytes(data)
	if err != nil {
		return errors.New(fmt.Sprintf("%s, %+v", err.Error(), valErr))
	}
	return nil
}

// Validate given yaml data against loaded jsonschema.
func (j *Schema) ValidateYAML(data []byte) error {
	json_data, err := yaml.YAMLToJSON(data)
	if err != nil {
		return err
	}
	return j.ValidateJSON(json_data)
}
