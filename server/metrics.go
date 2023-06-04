package server

import "github.com/prometheus/client_golang/prometheus"

// Metrics contsins all metric types
type Metrics struct {
	ReuestCount  prometheus.Counter
	ReuestErrors prometheus.Counter
}

func RegisterPrometheusMetrics() Metrics {

	ReuestCount := prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "server_requests_count",
			Help: "The total number of requests",
		},
	)

	ReuestErrors := prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "server_requests_errors",
			Help: "The total number of failed requests",
		},
	)

	prometheus.MustRegister(ReuestCount)
	prometheus.MustRegister(ReuestErrors)

	return Metrics{ReuestCount: ReuestCount, ReuestErrors: ReuestErrors}
}
