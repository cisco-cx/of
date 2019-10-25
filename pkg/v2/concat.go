package v2

import "io"

// Interface to merge all configs in conf.d dir.
type Concatenate interface {
	Concat() (io.Reader, error) // Join all configs in given Dir.
}
