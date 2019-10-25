package snmp

// Helps in identifing configs applicable for a OID.
type Lookup interface {
	Build() error                  // Build lookup map
	Find(string) ([]string, error) // For given OID, return array of configs applicable.
}
