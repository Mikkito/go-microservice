package metrics

import "github.com/prometheus/client_golang/prometheus"

var RateLimitExceeded = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "rate_limit_exceeded_total",
		Help: "Total number of rejected requests due to rate limiting",
	},
)
