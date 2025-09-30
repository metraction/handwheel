package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// Counter for total number of images processed
	ImagesProcessedTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "handwheel_images_processed_total",
		Help: "The total number of container images processed",
	})

	// Counter for total number of deployments posted to DevLake
	DeploymentsPostedTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "handwheel_deployments_posted_total",
		Help: "The total number of deployments posted to DevLake",
	})

	// Counter for Prometheus query executions
	PrometheusQueriesTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "handwheel_prometheus_queries_total",
		Help: "The total number of Prometheus queries executed",
	}, []string{"status"})

	// Gauge for number of unique images seen
	UniqueImagesGauge = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "handwheel_unique_images_current",
		Help: "The current number of unique images being tracked",
	})

	// Histogram for processing duration
	ProcessingDurationSeconds = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "handwheel_processing_duration_seconds",
		Help:    "Time taken to process image metrics pipeline",
		Buckets: prometheus.DefBuckets,
	})

	// Counter for errors by type
	ErrorsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "handwheel_errors_total",
		Help: "The total number of errors by type",
	}, []string{"type", "component"})
)
