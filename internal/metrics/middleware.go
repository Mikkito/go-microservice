package metrics

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		recorder := &statusRecorder{
			ResponseWriter: w,
			status:         http.StatusOK,
		}

		next.ServeHTTP(recorder, r)

		duration := time.Since(start).Seconds()
		status := recorder.status

		RequestsTotal.WithLabelValues(
			r.Method,
			r.URL.Path,
			http.StatusText(status),
		).Inc()

		RequestDuration.WithLabelValues(
			r.Method,
			r.URL.Path,
		).Observe(duration)

		if status >= 400 {
			ErrorsTotal.WithLabelValues(
				r.Method,
				r.URL.Path,
				http.StatusText(status),
			).Inc()
		}
	})
}

func Handler() http.Handler {
	return promhttp.Handler()
}
