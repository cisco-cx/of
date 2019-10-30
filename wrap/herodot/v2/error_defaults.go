package v2

// Package `github.com/cisco-cx/of/wrap/herodo/v2` is the
// `v2` named-version wrapper of the [herodot](https://"github.com/ory/herodot")
// standard library package.

import (
	"github.com/ory/herodot"
	http "github.com/cisco-cx/of/wrap/http/v2"
)

var ErrNotFound = herodot.DefaultError{
	StatusField: http.StatusText(http.StatusNotFound),
	ErrorField:  "The requested resource could not be found",
	CodeField:   http.StatusNotFound,
}

var ErrUnauthorized = herodot.DefaultError{
	StatusField: http.StatusText(http.StatusUnauthorized),
	ErrorField:  "The request could not be authorized",
	CodeField:   http.StatusUnauthorized,
}

var ErrForbidden = herodot.DefaultError{
	StatusField: http.StatusText(http.StatusForbidden),
	ErrorField:  "The requested action was forbidden",
	CodeField:   http.StatusForbidden,
}

var ErrInternalServerError = herodot.DefaultError{
	StatusField: http.StatusText(http.StatusInternalServerError),
	ErrorField:  "An internal server error occurred, please contact the system administrator",
	CodeField:   http.StatusInternalServerError,
}

var ErrBadRequest = herodot.DefaultError{
	StatusField: http.StatusText(http.StatusBadRequest),
	ErrorField:  "The request was malformed or contained invalid parameters",
	CodeField:   http.StatusBadRequest,
}

var ErrUnsupportedMediaType = herodot.DefaultError{
	StatusField: http.StatusText(http.StatusUnsupportedMediaType),
	ErrorField:  "The request is using an unknown content type",
	CodeField:   http.StatusUnsupportedMediaType,
}

var ErrConflict = herodot.DefaultError{
	StatusField: http.StatusText(http.StatusConflict),
	ErrorField:  "The resource could not be created due to a conflict",
	CodeField:   http.StatusConflict,
}
