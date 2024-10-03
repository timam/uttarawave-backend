package metrics

import (
	"expvar"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"runtime"
)

var (
	httpRequestsTotal      *prometheus.CounterVec
	httpRequestDuration    *prometheus.HistogramVec
	httpRequestSize        *prometheus.HistogramVec
	httpResponseSize       *prometheus.HistogramVec
	httpRequestRate        *prometheus.CounterVec
	httpErrorRate          *prometheus.CounterVec
	httpRequestQueueLength prometheus.Gauge
	httpConnectionCount    prometheus.Gauge
)

func InitializeMetrics() error {
	// Initialize Prometheus metrics
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request latencies in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path", "status"},
	)

	httpRequestSize = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_size_bytes",
			Help:    "Size of HTTP requests in bytes",
			Buckets: prometheus.ExponentialBuckets(100, 10, 8),
		},
		[]string{"method", "path"},
	)

	httpResponseSize = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_response_size_bytes",
			Help:    "Size of HTTP responses in bytes",
			Buckets: prometheus.ExponentialBuckets(100, 10, 8),
		},
		[]string{"method", "path", "status"},
	)

	httpRequestRate = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_rate",
			Help: "Number of HTTP requests per second",
		},
		[]string{"method", "path"},
	)

	httpErrorRate = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_error_rate",
			Help: "Number of HTTP requests resulting in errors",
		},
		[]string{"method", "path", "status"},
	)

	httpRequestQueueLength = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "http_request_queue_length",
			Help: "Number of HTTP requests waiting to be processed",
		},
	)

	httpConnectionCount = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "http_connection_count",
			Help: "Number of active HTTP connections",
		},
	)

	prometheus.MustRegister(httpRequestsTotal, httpRequestDuration, httpRequestSize, httpResponseSize, httpRequestRate, httpErrorRate, httpRequestQueueLength, httpConnectionCount)

	// Register Go runtime metrics
	expvar.Publish("goroutines", expvar.Func(func() interface{} {
		return runtime.NumGoroutine()
	}))
	expvar.Publish("memory", expvar.Func(func() interface{} {
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		return m.Alloc
	}))

	return nil
}

func RecordHTTPMetrics(method, path, status string, duration, requestSize, responseSize float64) {
	httpRequestsTotal.WithLabelValues(method, path, status).Inc()
	httpRequestDuration.WithLabelValues(method, path, status).Observe(duration)
	httpRequestSize.WithLabelValues(method, path).Observe(requestSize)
	httpResponseSize.WithLabelValues(method, path, status).Observe(responseSize)
	httpRequestRate.WithLabelValues(method, path).Inc()

	if status[0] == '4' || status[0] == '5' {
		httpErrorRate.WithLabelValues(method, path, status).Inc()
	}
}

func MetricsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		promhttp.Handler().ServeHTTP(c.Writer, c.Request)
	}
}
