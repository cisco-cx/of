package snmp

import (
	of "github.com/cisco-cx/of/pkg/v2"
)

// Helps in identifing configs applicable for a OID.
type Lookup interface {
	Build() error                         // Build lookup map
	Find(*[]of.TrapVar) ([]string, error) // For given OID, return array of configs applicable.
}
