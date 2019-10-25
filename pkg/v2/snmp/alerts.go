package snmp

import (
	of "github.com/cisco-cx/of/pkg/v2"
)

// Parse configs and generate alerts for give OID
type AlertGenerator interface {
	Alert(string, []string) of.Alert // Generate v1.Alert for given OID and array of Config names.
}
