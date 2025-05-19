package integrations

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/Tiktai/handler/model"
)

type PrometheusIntegration struct {
	config model.PrometheusConfig
}

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

func PrometheusResult(o model.OutputWithError) []model.ImageMetric {
	if o.Err != nil {
		log.Println("Error fetching metrics:", o.Err)
		return nil
	} else {
		return o.Result.([]model.ImageMetric)
	}
}

func NewPrometheusIntegration(cfg *model.Config) *PrometheusIntegration {
	return &PrometheusIntegration{
		config: cfg.Prometheus,
	}
}

func (integration PrometheusIntegration) FetchImageMetrics(_ any) model.OutputWithError {
	log.Println("Fetching metrics.")
	// Example: http://localhost:9090/api/v1/query?query=kube_pod_container_info
	url := fmt.Sprintf("%s/api/v1/query?query=kube_pod_container_info", integration.config.URL)
	resp, err := http.Get(url)
	if err != nil {
		return model.OutputWithError{Err: err}
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return model.OutputWithError{Err: err}
	}
	var promResp PrometheusResponse
	if err := json.Unmarshal(body, &promResp); err != nil {
		return model.OutputWithError{Err: err}
	}
	var metrics []model.ImageMetric
	for _, r := range promResp.Data.Result {
		image_spec := r.Metric["image_spec"]
		if image_spec == "" {
			image_spec = r.Metric["pod"] // fallback, adjust as needed
		}
		if image_spec != "" {
			metrics = append(metrics, model.ImageMetric{Image_spec: image_spec})
		}
	}
	return model.OutputWithError{Result: metrics}
}
