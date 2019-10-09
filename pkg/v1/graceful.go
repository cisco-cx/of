package v1

// Methods required to implement graceful.
type Graceful interface {
	Start() error
	Stop() error
}
