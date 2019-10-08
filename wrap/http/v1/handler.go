package v1

import (
	"net/http"

	of "github.com/cisco-cx/of/lib/v1"
)

// Interface to convert of.* to http.*
type handlerOveride struct {
	serverHTTP func(of.ResponseWriter, of.Request)
}

// Pass on http.ResponseWriter and http.Request to of.Handler
func (h *handlerOveride) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	h.serverHTTP(rw, r)
}
