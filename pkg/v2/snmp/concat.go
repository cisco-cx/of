package snmp

type Concatenate interface {
	Concat(string) V2Config // Join all configs in given Dir.
}
