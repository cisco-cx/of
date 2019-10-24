package snmp

import (
	v1 "github.com/cisco-cx/of/pkg/v1"
)

type AlertGenerator interface {
	Alert(string, []string) v1.Alert // Generate v1.Alert for given OID and array of Config names.
}
