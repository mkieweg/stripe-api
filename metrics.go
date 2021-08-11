package stripeapi

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	registerHttpOK = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "frontend api",
		Subsystem: "register",
		Name:      "http_ok_total",
	})

	subscribeHttpOK = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "frontend api",
		Subsystem: "subscribe",
		Name:      "http_ok_total",
	})

	webhookHttpOK = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "frontend api",
		Subsystem: "webhook",
		Name:      "http_ok_total",
	})

	registerHttpInternalServerError = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "frontend api",
		Subsystem: "register",
		Name:      "http_internal_server_error_total",
	})

	subscribeHttpInternalServerError = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "frontend api",
		Subsystem: "subscribe",
		Name:      "http_internal_server_error_total",
	})

	webhookHttpInternalServerError = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "frontend api",
		Subsystem: "webhook",
		Name:      "http_internal_server_error_total",
	})

	inFlightGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "in_flight_requests",
		Help: "A gauge of requests currently being served by the wrapped handler.",
	})

	counter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "api_requests_total",
			Help: "A counter for requests to the wrapped handler.",
		},
		[]string{"code", "method"},
	)

	duration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "request_duration_seconds",
			Help:    "A histogram of latencies for requests.",
			Buckets: []float64{.25, .5, 1, 2.5, 5, 10},
		},
		[]string{"handler", "method"},
	)

	responseSize = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "response_size_bytes",
			Help:    "A histogram of response sizes for requests.",
			Buckets: []float64{200, 500, 900, 1500},
		},
		[]string{},
	)

	registerChain = promhttp.InstrumentHandlerInFlight(inFlightGauge,
		promhttp.InstrumentHandlerDuration(duration.MustCurryWith(prometheus.Labels{"handler": "register"}),
			promhttp.InstrumentHandlerCounter(counter,
				promhttp.InstrumentHandlerResponseSize(responseSize, registerHandler),
			),
		),
	)

	subscribeChain = promhttp.InstrumentHandlerInFlight(inFlightGauge,
		promhttp.InstrumentHandlerDuration(duration.MustCurryWith(prometheus.Labels{"handler": "subscribe"}),
			promhttp.InstrumentHandlerCounter(counter,
				promhttp.InstrumentHandlerResponseSize(responseSize, subscribeHandler),
			),
		),
	)

	webhookChain = promhttp.InstrumentHandlerInFlight(inFlightGauge,
		promhttp.InstrumentHandlerDuration(duration.MustCurryWith(prometheus.Labels{"handler": "webhook"}),
			promhttp.InstrumentHandlerCounter(counter,
				promhttp.InstrumentHandlerResponseSize(responseSize, webhookHandler),
			),
		),
	)
)
