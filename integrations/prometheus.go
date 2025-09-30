package integrations

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/metraction/handwheel/metrics"
	"github.com/metraction/handwheel/model"
)

type PrometheusIntegration struct {
	config     model.PrometheusConfig
	httpClient *http.Client
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

// normalizePEM reformats a PEM string to ensure proper newlines after header, before footer, and every 64 chars in base64 body.
// Added for argocd pecularities
func normalizePEM(pemBytes []byte) []byte {
	pemStr := string(pemBytes)
	header := "-----BEGIN CERTIFICATE-----"
	footer := "-----END CERTIFICATE-----"
	pemStr = strings.TrimSpace(pemStr)
	if strings.Contains(pemStr, header) && strings.Contains(pemStr, footer) {
		headIdx := strings.Index(pemStr, header)
		footIdx := strings.Index(pemStr, footer)
		if headIdx != -1 && footIdx != -1 {
			// Extract header, base64, footer
			head := header
			foot := footer
			base := pemStr[headIdx+len(header) : footIdx]
			base = strings.ReplaceAll(base, "\n", " ")
			base = strings.ReplaceAll(base, "\r", " ")
			base = strings.ReplaceAll(base, " ", "")
			// Insert newlines every 64 chars
			var lines []string
			for i := 0; i < len(base); i += 64 {
				end := i + 64
				if end > len(base) {
					end = len(base)
				}
				lines = append(lines, base[i:end])
			}
			return []byte(head + "\n" + strings.Join(lines, "\n") + "\n" + foot + "\n")
		}
	}
	return pemBytes
}

func NewPrometheusIntegration(cfg *model.Config) *PrometheusIntegration {
	transport := NewHttpTransport(cfg)
	client := &http.Client{
		Transport: transport,
	}

	return &PrometheusIntegration{
		config:     cfg.Prometheus,
		httpClient: client,
	}
}

func (integration PrometheusIntegration) FetchImageMetrics(_ any) model.OutputWithError {
	log.Println("Fetching metrics.")
	if integration.config.Query == "" {
		integration.config.Query = "kube_pod_container_info" // Default query if not specified
	}
	url := fmt.Sprintf("%s/api/v1/query?query=%s", integration.config.URL, integration.config.Query)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		metrics.PrometheusQueriesTotal.WithLabelValues("error").Inc()
		metrics.ErrorsTotal.WithLabelValues("http_request", "prometheus").Inc()
		return model.OutputWithError{Err: err}
	}
	if integration.config.Auth.Token == "" && integration.config.Auth.Username != "" {
		log.Println("Using token and username for authentication.")
		req.SetBasicAuth(integration.config.Auth.Username, integration.config.Auth.Password)
	} else if integration.config.Auth.Token != "" {
		log.Println("Using token for authentication.")
		req.Header.Set("Authorization", "Bearer "+integration.config.Auth.Token)
	} else {
		log.Println("No authentication configured, using anonymous access.")
	}
	resp, err := integration.httpClient.Do(req)
	if err != nil {
		metrics.PrometheusQueriesTotal.WithLabelValues("error").Inc()
		metrics.ErrorsTotal.WithLabelValues("http_do", "prometheus").Inc()
		return model.OutputWithError{Err: err}
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		metrics.PrometheusQueriesTotal.WithLabelValues("error").Inc()
		metrics.ErrorsTotal.WithLabelValues("read_body", "prometheus").Inc()
		return model.OutputWithError{Err: err}
	}
	var promResp PrometheusResponse
	if err := json.Unmarshal(body, &promResp); err != nil {
		metrics.PrometheusQueriesTotal.WithLabelValues("error").Inc()
		metrics.ErrorsTotal.WithLabelValues("json_unmarshal", "prometheus").Inc()
		return model.OutputWithError{Err: err}
	}
	var imageMetrics []model.ImageMetric
	for _, r := range promResp.Data.Result {
		image_spec := r.Metric["image_spec"]
		if image_spec == "" {
			image_spec = r.Metric["pod"] // fallback, adjust as needed
		}
		labels := make(map[string]string)
		for k, v := range r.Metric {
			labels[k] = v
		}
		if image_spec != "" {
			imageMetrics = append(imageMetrics, model.ImageMetric{Image_spec: image_spec, Labels: labels})
		}
	}
	
	// Record successful query and number of images processed
	metrics.PrometheusQueriesTotal.WithLabelValues("success").Inc()
	metrics.ImagesProcessedTotal.Add(float64(len(imageMetrics)))
	
	return model.OutputWithError{Result: imageMetrics}
}
