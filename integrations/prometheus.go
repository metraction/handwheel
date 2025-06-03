package integrations

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

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
	client := &http.Client{Transport: transport}

	return &PrometheusIntegration{
		config:     cfg.Prometheus,
		httpClient: client,
	}
}

func (integration PrometheusIntegration) FetchImageMetrics(_ any) model.OutputWithError {
	log.Println("Fetching metrics.")
	url := fmt.Sprintf("%s/api/v1/query?query=kube_pod_container_info", integration.config.URL)

	resp, err := integration.httpClient.Get(url)
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
