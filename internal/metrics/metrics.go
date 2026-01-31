package metrics

import "github.com/prometheus/client_golang/prometheus"

var HttpRequestsTotal = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Total number of HTTP requests",
	},
	[]string{"path","method", "status"},
)
var HttpRequestDuration = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "http_request_duration_seconds",
		Help:    "Duration of HTTP requests",
		Buckets: prometheus.DefBuckets,
	},
	[]string{"path", "method"},
)

var RateLimitedRequestsTotal = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "rate_limited_requests_total",
		Help: "Total number of HTTP requests rejected due to rate limiting",
	},
	[]string{"path"},
)

var GatewayErrorsTotal = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "gateway_errors_total",
		Help: "Total number of errors returned by the gateway",
	},
	[]string{"path", "status"},
)

var RetryAttemptsTotal = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "gateway_retry_attempts_total",
		Help: "Total number of retry attempts made by the gateway",
	},
	[]string{"path"},
)




func Init() {
	prometheus.MustRegister(HttpRequestsTotal)
	prometheus.MustRegister(HttpRequestDuration)
	prometheus.MustRegister(RateLimitedRequestsTotal)
	prometheus.MustRegister(GatewayErrorsTotal)
	prometheus.MustRegister(RetryAttemptsTotal)
}

