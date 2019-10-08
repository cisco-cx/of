package v1

// Counter errors.
const (
	ErrCounterCreateFailed  = Error("Failed to create counter.")
	ErrCounterDestroyFailed = Error("Failed to remove counter.")
)

// Error represents an OF error.
type Error string

// Error returns the error as a string.
func (e Error) Error() string { return string(e) }
