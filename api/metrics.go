package api

import (
	"net/http"

	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// Mutex for thread-safe access to metrics
	mutex = &sync.Mutex{}

	// Prometheus metrics
	totalRequests = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "ars2024_project_total_requests",
			Help: "Total number of requests for the past 24 hours.",
		},
	)
	successfulRequests = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "ars2024_project_successful_requests",
			Help: "Total number of successful http hits for the past 24 hours.",
		},
	)
	failedRequests = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "ars2024_project_failed_requests",
			Help: "Total number of failed http hits for the past 24 hours.",
		},
	)
	requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "ars2024_project_request_duration_seconds",
			Help:    "Histogram of response latency (seconds) of http requests.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)
	requestRate = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "ars2024_project_request_rate",
			Help: "Request rate per endpoint for the past 24 hours.",
		},
		[]string{"method", "endpoint"},
	)

	// Map to store request timestamps for request rate calculation
	requestTimestamps = make(map[string]time.Time)
)

func init() {
	// Register metrics that will be exposed.
	prometheus.MustRegister(totalRequests, successfulRequests, failedRequests, requestDuration, requestRate)
}

// MetricsHandler returns an HTTP handler for serving Prometheus metrics.
func MetricsHandler() http.Handler {
	return promhttp.Handler()
}

// Count is a middleware function to count HTTP requests, measure request duration,
// and calculate request rate.
func Count(f http.HandlerFunc, method, endpoint string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Increment total request count
		mutex.Lock()
		totalRequests.Inc()
		mutex.Unlock()

		// Execute the handler
		rr := &responseRecorder{w, http.StatusOK}
		f(rr, r)

		// Measure request duration
		duration := time.Since(start).Seconds()
		requestDuration.WithLabelValues(method, endpoint).Observe(duration)

		// Increment successful or failed request count based on status code
		if statusCode := rr.statusCode; statusCode >= 200 && statusCode < 400 {
			mutex.Lock()
			successfulRequests.Inc()
			mutex.Unlock()
		} else {
			mutex.Lock()
			failedRequests.Inc()
			mutex.Unlock()
		}

		// Update request rate
		updateRequestRate(method, endpoint)
	}
}

// updateRequestRate updates the request rate metric for a given method and endpoint.
func updateRequestRate(method, endpoint string) {
	mutex.Lock()
	defer mutex.Unlock()

	// Store current timestamp for the endpoint
	key := method + "_" + endpoint
	requestTimestamps[key] = time.Now()

	// Calculate request rate for the past 24 hours
	var count float64
	currentTime := time.Now()
	for _, timestamp := range requestTimestamps {
		if currentTime.Sub(timestamp) <= 24*time.Hour {
			count++
		}
	}
	requestRate.WithLabelValues(method, endpoint).Set(count)
}

// responseRecorder is a custom http.ResponseWriter to track the status code.
type responseRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (rec *responseRecorder) WriteHeader(code int) {
	rec.statusCode = code
	rec.ResponseWriter.WriteHeader(code)
}
