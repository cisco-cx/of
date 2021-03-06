package v2

import "io"

// Cache data
type Cacher interface {
	Read(io.Reader, interface{}) error
	Write(io.Writer, interface{}) error
}
