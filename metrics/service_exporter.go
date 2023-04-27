package metrics

import (
	"fmt"
	log "github.com/Heqiaomu/glog"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	_ "net/http/pprof"

	"github.com/prometheus/client_golang/prometheus"
)

// handler wraps an unfiltered http.Handler but uses a filtered handler,
// created on the fly, if filtering is requested. Create instances with
// newHandler.
type handler struct {
	metricsHandler http.Handler
	// exporterMetricsRegistry is a separate registry for the metrics about
	// the exporter itself.
	exporterMetricsRegistry     *prometheus.Registry
	maxRequests                 int
	serviceGoCollectorOpts      ServiceGoCollectorOpts
	serviceProcessCollectorOpts ServiceProcessCollectorOpts
}
type ServiceCollectorOpts struct {
	LabelsKey   []string
	LabelsValue []string
}

func NewHandler(maxRequests int, opts ServiceCollectorOpts) *handler {
	h := &handler{
		exporterMetricsRegistry:     prometheus.NewRegistry(),
		maxRequests:                 maxRequests,
		serviceGoCollectorOpts:      ServiceGoCollectorOpts{},
		serviceProcessCollectorOpts: ServiceProcessCollectorOpts{},
	}
	h.serviceGoCollectorOpts.LabelsKey = opts.LabelsKey
	h.serviceGoCollectorOpts.LabelsValue = opts.LabelsValue
	h.serviceProcessCollectorOpts.LabelsKey = opts.LabelsKey
	h.serviceProcessCollectorOpts.LabelsValue = opts.LabelsValue
	return h
}

// ServeHTTP implements http.Handler.
func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler, err := h.innerHandler()
	if err != nil {
		log.Errorf("Manage err when creating http handler. Now to return 400. Err: [%v]", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("Couldn't create filtered metrics handler: %s", err)))
		return
	}
	handler.ServeHTTP(w, r)
}

// innerHandler is used to create both the one unfiltered http.Handler to be
// wrapped by the outer handler and also the filtered handlers created on the
// fly. The former is accomplished by calling innerHandler without any arguments
// (in which case it will log all the collectors enabled via command-line
// flags).
func (h *handler) innerHandler() (http.Handler, error) {
	r := prometheus.NewRegistry()
	if err := r.Register(NewServiceGoCollector(h.serviceGoCollectorOpts)); err != nil {
		return nil, fmt.Errorf("couldn't register service go collector: %s", err)
	}
	if err := r.Register(NewServiceProcessCollector(h.serviceProcessCollectorOpts)); err != nil {
		return nil, fmt.Errorf("couldn't register service process collector: %s", err)
	}
	handler := promhttp.HandlerFor(
		prometheus.Gatherers{h.exporterMetricsRegistry, r},
		promhttp.HandlerOpts{
			ErrorHandling:       promhttp.ContinueOnError,
			MaxRequestsInFlight: h.maxRequests,
			Registry:            h.exporterMetricsRegistry,
		},
	)
	handler = promhttp.InstrumentMetricHandler(
		h.exporterMetricsRegistry, handler)
	return handler, nil

}
