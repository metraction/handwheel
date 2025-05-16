package integrations

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/Tiktai/handler/model"
	"github.com/Tiktai/handler/streams"
)

// PrometheusResponse is a minimal struct for parsing Prometheus API responses
type PrometheusResponse struct {
	Status string `json:"status"`
	Data   struct {
		ResultType string `json:"resultType"`
		Result     []struct {
			Metric map[string]string `json:"metric"`
			// Value: [ <timestamp>, "<value>" ]
			Value []interface{} `json:"value"`
		} `json:"result"`
	} `json:"data"`
}

// FetchAndSubmitPeriodically fetches kube_pod_container_info and submits ImageMetrics to the stream
func FetchAndSubmitPeriodically(cfg *model.Config) {
	promURL := cfg.Prometheus.URL
	interval, err := time.ParseDuration(cfg.Prometheus.Interval)
	if err != nil {
		log.Printf("Invalid prometheus.interval in config: %v, defaulting to 1m", err)
		interval = time.Minute
	}
	batchSize := cfg.Prometheus.BatchSize
	if batchSize <= 0 {
		log.Printf("Invalid prometheus.batch_size in config: %d, defaulting to 10", batchSize)
		batchSize = 10
	}
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		log.Println("Fetching metrics.")
		<-ticker.C
		metrics, err := fetchImageMetrics(promURL)
		if err != nil {
			log.Printf("Error fetching metrics: %v", err)
			continue
		}
		// Batch submission, adjust as needed
		for i := 0; i < len(metrics); i += batchSize {
			end := i + batchSize
			if end > len(metrics) {
				end = len(metrics)
			}
			batch := metrics[i:end]
			// Submit to your go-stream batch processor
			go func(batch []streams.ImageMetric) {
				_ = streams.ProcessImage(batch) // or AugmentImages(batch) or your custom logic
			}(batch)
		}
	}
}

func fetchImageMetrics(promURL string) ([]streams.ImageMetric, error) {
	// Example: http://localhost:9090/api/v1/query?query=kube_pod_container_info
	url := fmt.Sprintf("%s/api/v1/query?query=kube_pod_container_info", promURL)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var promResp PrometheusResponse
	if err := json.Unmarshal(body, &promResp); err != nil {
		return nil, err
	}
	var metrics []streams.ImageMetric
	for _, r := range promResp.Data.Result {
		sha := r.Metric["container_id"]
		if sha == "" {
			sha = r.Metric["pod"] // fallback, adjust as needed
		}
		if sha != "" {
			metrics = append(metrics, streams.ImageMetric{Sha: sha})
		}
	}
	return metrics, nil
}
