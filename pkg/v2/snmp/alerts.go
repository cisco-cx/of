package snmp

import (
	of "github.com/cisco-cx/of/pkg/v2"
)

// Parse given configs and generate alerts.
type AlertGenerator interface {
	Alert([]string) []of.Alert // Generate an array of of.Alert for given array of Config names.
}
