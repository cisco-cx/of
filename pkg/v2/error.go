package v2

const (

	// Value errors.
	ErrOIDNotFound      = Error("OID not present in trap vars.")
	ErrUnknownAs        = Error("Unknown v2.snmp.As type.")
	ErrNoneNumericalOID = Error("Numerical OID expected..")

	// Concatenate errors.
	ErrPathIsNotDir = Error("Path is not a directory.")

	// Mod errors.
	ErrInvalidOperation = Error("Operation not possible for given Mod.")
	ErrKeyMissing       = Error("Key missing in mod.")

	// Alert Generator errors.
	ErrConfigNotFound   = Error("Unknown config.")
	ErrNoMatch          = Error("No alert matched in alert config.")
	ErrUnknownEventType = Error("Unknown event type specified.")

	// Counter errors.
	ErrCounterCreateFailed  = Error("Failed to create counter.")
	ErrCounterDestroyFailed = Error("Failed to remove counter.")
)

// Error represents an OF error.
type Error string

// Error returns the error as a string.
func (e Error) Error() string { return string(e) }
