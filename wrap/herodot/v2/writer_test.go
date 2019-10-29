package v2_test

import (
	"testing"

	of "github.com/cisco-cx/of/pkg/v2"
	hero "github.com/cisco-cx/of/wrap/herodot/v2"
)

// Enforce interface
func TestInterface(t *testing.T) {
	var _ of.Writer = &hero.Writer{}
}
