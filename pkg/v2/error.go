package v2

// Value errors.
const (
	ErrOIDNotFound      = Error("OID not present in trap vars.")
	ErrUnknownAs        = Error("Unknown v2.snmp.As type.")
	ErrNoneNumericalOID = Error("Numerical OID expected..")
)

// Error represents an OF error.
type Error string

// Error returns the error as a string.
func (e Error) Error() string { return string(e) }
