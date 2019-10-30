package v2

import (
	"fmt"
	"net/url"

	of "github.com/cisco-cx/of/pkg/v2"
	health "github.com/cisco-cx/of/wrap/health/v1"
	http "github.com/cisco-cx/of/wrap/http/v2"
	logger "github.com/cisco-cx/of/wrap/logrus/v2"
	prometheus "github.com/cisco-cx/of/wrap/prometheus/client_golang/v2"
)

type Handler struct {
	Config *of.SNMPConfig
	server *http.Server
	SNMP   *Service
	Log    *logger.Logger
}

func (h *Handler) Run() {

	httpConfig := of.HTTPConfig{
		ListenAddress: h.Config.ListenAddress,
		ReadTimeout:   h.Config.AMTimeout,
		WriteTimeout:  h.Config.AMTimeout,
	}

	// Add health check.
	hc := health.New()
	amUrl, err := url.Parse(h.Config.AMAddress)
	if err != nil {
		h.Log.WithError(err).Fatalf("Failed to parse the Alerts manager address: %s", h.Config.AMAddress)
	}
	healthUrl, _ := url.Parse("/-/healthy")
	err = hc.AddURL("health_check", amUrl.ResolveReference(healthUrl).String(), h.Config.AMTimeout)
	if err != nil {
		h.Log.WithError(err).Fatalf("Failed to add health check.")
	}

	// Configure HTTP server to handle various requests.
	h.server = http.NewServer(&httpConfig)

	h.server.HandleFunc("/", func(w of.ResponseWriter, r of.Request) {
		fmt.Fprint(w, h.Config.Version)
	})

	// Handling health check.
	h.server.HandleFunc("/health", func(w of.ResponseWriter, r of.Request) {
		err := hc.State("health_check")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, err.Error())
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	// Handling prometheus calls.
	h.server.Handle("/metrics", prometheus.NewHandler())

	// Handling status calls.
	h.server.HandleFunc("/api/v2/status", func(w of.ResponseWriter, r of.Request) {
		fmt.Fprint(w, of.AppStatus{
			ApiVersion:  "",
			Description: "AlertManager Client for SNMP Traps",
			Links:       of.AppStatusLinks{About: "https://github.com/cisco-cx/am-client-snmp"},
			Status:      "success",
		})
	})

	h.server.HandleFunc("/api/v2/events", h.SNMP.AlertHandler)

	// Starting health check.
	err = hc.Start()
	if err != nil {
		h.Log.WithError(err).Fatalf("Failed to start at health check.")
	}
	defer hc.Stop()

	// Starting SNMP server.
	err = h.server.ListenAndServe()
	if err != nil {
		h.Log.WithError(err).Fatalf("Failed to listen at %s", h.Config.ListenAddress)
	}

}
