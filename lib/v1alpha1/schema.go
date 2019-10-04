package v1alpha1

// Defines schema loader
type SchemaLoader interface {
	Load([]byte) error
}

// Defines schema Validator
type SchemaValidator interface {
	ValidateYAML([]byte) error
	ValidateJSON([]byte) error
}
