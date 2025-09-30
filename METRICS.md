# Prometheus Metrics

This application exposes Prometheus metrics on the `/metrics` endpoint of the health server.

## Endpoint

The metrics are available at: `http://localhost:<health_server_port>/metrics`

By default, the health server runs on port `8081`, so metrics would be available at:
`http://localhost:8081/metrics`

## Available Metrics

### Counters

- **`handwheel_images_processed_total`** - The total number of container images processed
- **`handwheel_deployments_posted_total`** - The total number of deployments posted to DevLake
- **`handwheel_prometheus_queries_total`** - The total number of Prometheus queries executed
  - Labels: `status` (success/error)
- **`handwheel_errors_total`** - The total number of errors by type
  - Labels: `type` (error type), `component` (component name)

### Gauges

- **`handwheel_unique_images_current`** - The current number of unique images being tracked

### Histograms

- **`handwheel_processing_duration_seconds`** - Time taken to process image metrics pipeline

## Error Types

The `handwheel_errors_total` metric tracks various error types:

### Prometheus Component
- `http_request` - Failed to create HTTP request
- `http_do` - HTTP request execution failed
- `read_body` - Failed to read response body
- `json_unmarshal` - Failed to parse JSON response

### DevLake Component
- `no_matching_project` - No matching DevLake project found for image
- `json_marshal` - Failed to marshal payload to JSON
- `http_request` - Failed to create HTTP request
- `http_do` - HTTP request execution failed
- `webhook_error` - DevLake webhook returned error status

## Usage with Prometheus

Add the following to your Prometheus configuration:

```yaml
scrape_configs:
  - job_name: 'handwheel'
    static_configs:
      - targets: ['localhost:8081']
    metrics_path: '/metrics'
    scrape_interval: 30s
```

## Sample Queries

```promql
# Rate of images processed per second
rate(handwheel_images_processed_total[5m])

# Error rate by component
rate(handwheel_errors_total[5m])

# Current number of unique images
handwheel_unique_images_current

# Prometheus query success rate
rate(handwheel_prometheus_queries_total{status="success"}[5m]) / rate(handwheel_prometheus_queries_total[5m])
```
