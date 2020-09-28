package v2_test

import (
	"testing"

	of "github.com/cisco-cx/of/pkg/v2"
	http "github.com/cisco-cx/of/wrap/http/v2"
)

// Enforce interface implementation.
func TestResponseWriterInterface(t *testing.T) {
	var _ of.ResponseWriter = &http.ResponseWriter{}
}
