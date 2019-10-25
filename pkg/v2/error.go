package v2

const (

	// Value errors.
	ErrOIDNotFound      = Error("OID not present in trap vars.")
	ErrUnknownAs        = Error("Unknown v2.snmp.As type.")
	ErrNoneNumericalOID = Error("Numerical OID expected..")

	// Concatenate errors.
	ErrPathIsNotDir = Error("Path is not a directory.")
)

// Error represents an OF error.
type Error string

// Error returns the error as a string.
func (e Error) Error() string { return string(e) }
