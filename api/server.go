package api

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Prometheus metrics
var (
	httpHits = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "ars2024_project_http_hit_total",
			Help: "Total number of http hits.",
		},
		[]string{"method", "endpoint"},
	)
	successfulHits = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "ars2024_project_http_success_total",
			Help: "Total number of successful http hits.",
		},
		[]string{"method", "endpoint"},
	)
	errorHits = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "ars2024_project_http_error_total",
			Help: "Total number of error http hits.",
		},
		[]string{"method", "endpoint"},
	)
	requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "ars2024_project_request_duration_seconds",
			Help:    "Histogram of response latency (seconds) of http requests.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)
	metricsList        = []prometheus.Collector{httpHits, successfulHits, errorHits, requestDuration}
	prometheusRegistry = prometheus.NewRegistry()
)

func init() {
	// Register metrics that will be exposed.
	prometheusRegistry.MustRegister(metricsList...)
}

func MetricsHandler() http.Handler {
	return promhttp.HandlerFor(prometheusRegistry, promhttp.HandlerOpts{})
}

func Count(f http.HandlerFunc, method, endpoint string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		timer := prometheus.NewTimer(requestDuration.WithLabelValues(method, endpoint))
		defer timer.ObserveDuration()

		httpHits.WithLabelValues(method, endpoint).Inc()

		rr := &responseRecorder{w, http.StatusOK}
		f(rr, r)

		statusCode := rr.statusCode
		if statusCode >= 200 && statusCode < 400 {
			successfulHits.WithLabelValues(method, endpoint).Inc()
		} else {
			errorHits.WithLabelValues(method, endpoint).Inc()
		}
	}
}

type responseRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (rec *responseRecorder) WriteHeader(code int) {
	rec.statusCode = code
	rec.ResponseWriter.WriteHeader(code)
}

func RunServer(router http.Handler) {
	server := &http.Server{
		Addr:    "0.0.0.0:8000",
		Handler: router,
	}

	log.Println("Server started on http://localhost:8000")
	log.Println("Swagger UI available at http://localhost:8000/swagger/index.html")

	// Start the server in a goroutine
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe error: %v", err)
		}
	}()
	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Server is shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	log.Println("Server exited gracefully")
}
