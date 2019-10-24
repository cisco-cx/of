package snmp

type Lookup interface {
	Build() error         // Build lookup map
	Find(string) []string // For given OID, return array of configs applicable.
}
