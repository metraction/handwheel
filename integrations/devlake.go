package integrations

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"time"

	"github.com/metraction/handwheel/metrics"
	"github.com/metraction/handwheel/model"
)

type DevLakeIntegration struct {
	config model.DevLakeConfig
	client *http.Client
}

func NewDevLakeIntegration(cfg *model.Config) *DevLakeIntegration {
	return &DevLakeIntegration{
		config: cfg.DevLake,
		client: &http.Client{
			Timeout:   10 * time.Second,
			Transport: NewHttpTransport(cfg),
		},
	}
}

func (di *DevLakeIntegration) PostDeployment(image model.Image) model.OutputWithError {
	repoURL := image.Labels["repo_url"]
	connectionID := -1
	for _, project := range di.config.Projects {
		for _, pattern := range project.Images {
			re, err := regexp.Compile(pattern)
			if err != nil {
				continue // skip invalid pattern
			}
			if re.MatchString(image.Image_spec) {
				connectionID = project.ConnectionID
				break
			}
		}
		if connectionID != -1 {
			break
		}
	}
	if connectionID == -1 {
		metrics.ErrorsTotal.WithLabelValues("no_matching_project", "devlake").Inc()
		return model.OutputWithError{Err: fmt.Errorf("no matching devlake project for repo_url: %s in %v", repoURL, di.config.Projects)}
	}

	url := fmt.Sprintf("%s/api/rest/plugins/webhook/connections/%d/deployments", di.config.URL, connectionID)
	token := di.config.Token

	// Prepare payload
	payload := map[string]any{
		"deploymentCommits": []map[string]string{
			{
				"commit_sha": image.Labels["commit_sha"], // expects commit_sha in labels
				"repo_url":   image.Labels["repo_url"],   // expects repo_url in labels
			},
		},
		"start_time": time.Now().UTC().Format(time.RFC3339), // or get from image.Labels if available
	}
	fmt.Println("Payload: ", payload)
	if v, ok := image.Labels["start_time"]; ok {
		payload["start_time"] = v
	}

	body, err := json.Marshal(payload)
	if err != nil {
		metrics.ErrorsTotal.WithLabelValues("json_marshal", "devlake").Inc()
		return model.OutputWithError{Err: fmt.Errorf("failed to marshal payload: %w", err)}
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		metrics.ErrorsTotal.WithLabelValues("http_request", "devlake").Inc()
		return model.OutputWithError{Err: fmt.Errorf("failed to create request: %w", err)}
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := di.client.Do(req)
	if err != nil {
		metrics.ErrorsTotal.WithLabelValues("http_do", "devlake").Inc()
		return model.OutputWithError{Err: fmt.Errorf("request failed: %w", err)}
	}
	defer resp.Body.Close()
	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 300 {
		metrics.ErrorsTotal.WithLabelValues("webhook_error", "devlake").Inc()
		return model.OutputWithError{Err: fmt.Errorf("devlake webhook error: %s", string(respBody))}
	}
	
	// Record successful deployment post
	metrics.DeploymentsPostedTotal.Inc()
	
	return model.OutputWithError{Result: string(respBody)}
}
