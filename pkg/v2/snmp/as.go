package snmp

// Interface to handle different types of MIB resolutions.
type ValueGenerator interface {
	ValueAs(string, As) (string, error)      // Compute value as `As` for given OID.
	Value(string) (string, error)            // Literal value for given OID.
	ValueStr(string) (string, error)         // String representation of the value, for given OID.
	ValueStrShort(string) (string, error)    // Short Name of the value, for given OID.
	OIDValue(string) (string, error)         // Literal value for OID pointed by given OID,
	OIDValueStr(string) (string, error)      // String representation of the value, for OID pointed by given OID.
	OIDValueStrShort(string) (string, error) // Short Name of the value, for OID pointed by given OID.
}
