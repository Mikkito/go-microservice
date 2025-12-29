package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	RequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP request",
		},
		[]string{"method", "path", "status"},
	)
	RequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "http_duration_seconds",
			Help: "Request duration",
		},
		[]string{"method", "path"},
	)
	ErrorsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_errors_total",
			Help: "Total number of HTTP errors",
		},
		[]string{"method", "path", "status"},
	)
)

func Init() {
	prometheus.MustRegister(
		RequestsTotal,
		RequestDuration,
		ErrorsTotal,
		RateLimitExceeded,
	)
}
