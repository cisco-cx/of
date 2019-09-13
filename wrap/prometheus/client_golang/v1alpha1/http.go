package v1alpha1

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	of "github.com/cisco-cx/of/lib/v1alpha1"
)

type Handler struct {
	ph http.Handler
}

func NewHandler() *Handler {
	return &Handler{
		ph: promhttp.Handler(),
	}
}

func (h *Handler) ServeHTTP(rw of.ResponseWriter, r of.Request) {
	h.ph.ServeHTTP(rw, r)
}
