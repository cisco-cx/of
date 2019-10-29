package v2

import (
	"github.com/ory/herodot"
	of "github.com/cisco-cx/of/pkg/v2"
	logger "github.com/cisco-cx/of/wrap/logrus/v2"
)

type Writer struct {
	h *herodot.JSONWriter
}

func New(l logger.Logger) *Writer {
	hero := herodot.NewJSONWriter(l.Logger())
	hero.ErrorEnhancer = nil
	return &Writer{h: hero}
}

func (wrt *Writer) Write(w of.ResponseWriter, r of.Request, e interface{}) {
	wrt.h.Write(w, r, e)
}

// WriteCode writes a response object to the ResponseWriter and sets a response code.
func (wrt *Writer) WriteCode(w of.ResponseWriter, r of.Request, code int, e interface{}) {
	wrt.h.WriteCode(w, r, code, e)
}

// WriteCreated writes a response object to the ResponseWriter with status code 201 and
// the Location header set to location.
func (wrt *Writer) WriteCreated(w of.ResponseWriter, r of.Request, location string, e interface{}) {
	wrt.h.WriteCreated(w, r, location, e)
}

// WriteError writes an error to ResponseWriter and tries to extract the error's status code by
// asserting statusCodeCarrier. If the error does not implement statusCodeCarrier, the status code
// is set to 500.
func (wrt *Writer) WriteError(w of.ResponseWriter, r of.Request, err interface{}) {
	wrt.h.WriteError(w, r, err)
}

// WriteErrorCode writes an error to ResponseWriter and forces an error code.
func (wrt *Writer) WriteErrorCode(w of.ResponseWriter, r of.Request, code int, err interface{}) {
	wrt.h.WriteErrorCode(w, r, code, err)
}
