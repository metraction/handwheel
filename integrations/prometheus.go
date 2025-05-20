package integrations

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/Tiktai/handler/model"
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

func NewPrometheusIntegration(cfg *model.Config) *PrometheusIntegration {
	pi := &PrometheusIntegration{
		config: cfg.Prometheus,
	}

	var caCertPool *x509.CertPool
	var err error
	if pi.config.CAFile != "" {
		log.Println("Using CA root file.", pi.config.CAFile)
		pem, err := os.ReadFile(pi.config.CAFile)
		if err != nil {
			log.Printf("failed to read CA root: %v", err)
			pi.httpClient = http.DefaultClient
			return pi
		}
		caCertPool, err = x509.SystemCertPool()
		if err != nil {
			log.Printf("failed to load system cert pool: %v", err)
			pi.httpClient = http.DefaultClient
			return pi
		}
		if !caCertPool.AppendCertsFromPEM(pem) {
			log.Printf("failed to append CA root cert")
			pi.httpClient = http.DefaultClient
			return pi
		}
	} else if pi.config.CARootPEM != "" {
		log.Println("Using CA root PEM.")
		caCertPool, err = x509.SystemCertPool()
		if err != nil {
			log.Printf("failed to load system cert pool: %v", err)
			pi.httpClient = http.DefaultClient
			return pi
		}
		if !caCertPool.AppendCertsFromPEM([]byte(pi.config.CARootPEM)) {
			log.Printf("failed to append CA root cert")
			pi.httpClient = http.DefaultClient
			return pi
		}
	}

	if caCertPool != nil {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{RootCAs: caCertPool},
		}
		pi.httpClient = &http.Client{Transport: tr}
	} else {
		pi.httpClient = http.DefaultClient
	}
	return pi
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
