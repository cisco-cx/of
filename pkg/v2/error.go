package v2

const (

	// Value errors.
	ErrOIDNotFound      = Error("OID not present in trap vars.")
	ErrUnknownAs        = Error("Unknown v2.snmp.As type.")
	ErrNoneNumericalOID = Error("Numerical OID expected..")

	// Mod errors.
	ErrInvalidOperation = Error("Operation not possible for given Mod.")
	ErrKeyMissing       = Error("Key missing in mod.")
)

// Error represents an OF error.
type Error string

// Error returns the error as a string.
func (e Error) Error() string { return string(e) }
