![deployments_frequency](https://handler-badges.enpace.ch/v1/Tiktai-handler/df)
[![lead_time_for_changes](https://handler-badges.enpace.ch/v1/Tiktai-handler/ltfc)](https://handler-badges.enpace.ch/v1/Tiktai-handler/ltfc-stats)

# Handler Project
[DevOps Research and Assessment (DORA)](https://cloud.google.com/blog/products/devops-sre/announcing-dora-2021-accelerate-state-of-devops-report) allows to measure how fast and reliable you can deliver value to your customers. Two of these metrics measure velocity:
- Deployment Frequency—How often an organization successfully releases to production
- Lead Time for Changes—The amount of time it takes a commit to get into production


Handler helps to get those metrics with least effort.
It tracks images reaching deployment environments and reports those events to devlake. 
This allows seamless reporting on DORA metrics without dependencies on pipelines.

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
