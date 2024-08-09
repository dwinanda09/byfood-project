package middleware

import (
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// HTTP request duration histogram
	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "http_request_duration_seconds",
			Help: "Duration of HTTP requests in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path", "status_code"},
	)

	// HTTP request counter
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status_code"},
	)

	// Active connections gauge
	activeConnections = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "http_active_connections",
			Help: "Number of active HTTP connections",
		},
	)

	// Database operation metrics
	databaseOperationDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "database_operation_duration_seconds",
			Help: "Duration of database operations in seconds",
			Buckets: []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10},
		},
		[]string{"operation", "table"},
	)

	databaseOperationsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "database_operations_total",
			Help: "Total number of database operations",
		},
		[]string{"operation", "table", "status"},
	)

	// Book-specific metrics
	booksTotal = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "books_total",
			Help: "Total number of books in the library",
		},
		[]string{"status"},
	)
)

// PrometheusMetrics middleware collects HTTP metrics
func PrometheusMetrics() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			
			// Increment active connections
			activeConnections.Inc()
			defer activeConnections.Dec()

			// Process request
			err := next(c)
			
			// Calculate request duration
			duration := time.Since(start).Seconds()
			
			// Get response status code
			statusCode := strconv.Itoa(c.Response().Status)
			
			// Record metrics
			httpRequestDuration.WithLabelValues(
				c.Request().Method,
				c.Path(),
				statusCode,
			).Observe(duration)
			
			httpRequestsTotal.WithLabelValues(
				c.Request().Method,
				c.Path(),
				statusCode,
			).Inc()

			return err
		}
	}
}

// RecordDatabaseOperation records database operation metrics
func RecordDatabaseOperation(operation, table string, duration time.Duration, success bool) {
	status := "success"
	if !success {
		status = "error"
	}
	
	databaseOperationDuration.WithLabelValues(operation, table).Observe(duration.Seconds())
	databaseOperationsTotal.WithLabelValues(operation, table, status).Inc()
}

// UpdateBookMetrics updates book-related metrics
func UpdateBookMetrics(totalBooks int, activeBooks int) {
	booksTotal.WithLabelValues("total").Set(float64(totalBooks))
	booksTotal.WithLabelValues("active").Set(float64(activeBooks))
}