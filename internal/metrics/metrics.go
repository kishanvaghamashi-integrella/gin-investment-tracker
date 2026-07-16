package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	// Counter: total number of HTTP requests
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "endpoint", "status"},
	)

	// Histogram: HTTP request duration in seconds
	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)

	// Gauge: current number of active requests
	httpRequestsActive = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "http_requests_active",
			Help: "Current number of active HTTP requests",
		},
	)
)

func Register() {
	prometheus.MustRegister(httpRequestsTotal)
	prometheus.MustRegister(httpRequestDuration)
	prometheus.MustRegister(httpRequestsActive)
}
