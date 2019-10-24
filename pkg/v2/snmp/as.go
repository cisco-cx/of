package snmp

type ValueGenerator interface {
	ValueAs(string, As) string      // Compute value as `As` for given OID.
	Value(string) string            // Literal value for given OID.
	ValueStr(string) string         // String representation of the value, for given OID.
	ValueStrShort(string) string    // Short Name of the value, for given OID.
	OIDValue(string) string         // Literal value for OID pointed by given OID,
	OIDValueStr(string) string      // String representation of the value, for OID pointed by given OID.
	OIDValueStrShort(string) string // Short Name of the value, for OID pointed by given OID.
}
