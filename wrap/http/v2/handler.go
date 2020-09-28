package v2

import (
	"net/http"

	of "github.com/cisco-cx/of/pkg/v2"
)

// Interface to convert of.* to http.*
type handlerOveride struct {
	m          of.Measurer
	serverHTTP func(of.ResponseWriter, of.Request)
}

// Pass on http.ResponseWriter and http.Request to of.Handler
func (h *handlerOveride) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	wrappedRW := NewResponseWriter(rw)
	h.m.Measure(wrappedRW, r, h.serverHTTP)
}
