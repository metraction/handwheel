![deployments_frequency](https://handler-badges.enpace.ch/v1/Tiktai-handler/df)
[![lead_time_for_changes](https://handler-badges.enpace.ch/v1/Tiktai-handler/ltfc)](https://handler-badges.enpace.ch/v1/Tiktai-handler/ltfc-stats)

# Handler Project
[DevOps Research and Assessment (DORA)](https://cloud.google.com/blog/products/devops-sre/announcing-dora-2021-accelerate-state-of-devops-report) enables organizations to measure the speed and reliability of customer value delivery. Two core velocity metrics include:

- Deployment Frequency — How frequently code deploys successfully to production.
- Lead Time for Changes — The duration from code commit to production deployment.

Handler simplifies capturing these metrics by automatically tracking container images across deployment environments and reporting events to DevLake. This eliminates manual pipeline instrumentation, enabling seamless DORA reporting with minimal configuration

<img src="docs/context.excalidraw.png" width="75%"/>

## Project Structure

- `cmd/` - Main application entry point(s)
  - `root.go` - Primary command logic
- `integrations/` - Integrations with external services
  - `crane.go` - Crane integration
  - `prometheus.go` - Prometheus metrics integration
- `logic/` - Core business logic (currently empty)
- `model/` - Data models and configuration
  - `config.go` - Configuration structures
  - `data.go` - Data types
- `routing/` - Routing logic
  - `devlake.go` - DevLake routing
  

## Getting Started

### Prerequisites
- Go 1.20 or newer

### Installation
1. Clone the repository:
   ```bash
   git clone <repository-url>
   cd handler
   ```
2. Install dependencies:
   ```bash
   go mod download
   ```

### Running Locally
To start the server:
```bash
go run main.go
```

## License
See [LICENSE](LICENSE) for details.

---
*Generated on 2025-05-17*
