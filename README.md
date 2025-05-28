![deployments_frequency](https://handler-badges.enpace.ch/v1/Tiktai-handler/df)
[![lead_time_for_changes](https://handler-badges.enpace.ch/v1/Tiktai-handler/ltfc)](https://handler-badges.enpace.ch/v1/Tiktai-handler/ltfc-stats)

# Handler Project
[DevOps Research and Assessment (DORA)](https://cloud.google.com/blog/products/devops-sre/announcing-dora-2021-accelerate-state-of-devops-report) enables organizations to measure the speed and reliability of customer value delivery. Two core velocity metrics include:

- Deployment Frequency — How frequently code deploys successfully to production.
- Lead Time for Changes — The duration from code commit to production deployment.

Handler simplifies capturing these metrics by automatically tracking container images across deployment environments and reporting events to DevLake. This eliminates manual pipeline instrumentation, enabling seamless DORA reporting with minimal configuration

<img src="docs/context.excalidraw.png" width="75%"/>

Check [dora-badges](https://github.com/Tiktai/dora-badge) for badges and stats.

# Getting Started

In order handler to work, docker images during build time need to be labeled, handler has to know where to fetch data from and where to push events to. 

## Label images 

During build time add following labels:

```
repo_url - where the docker is build
commit_sha - full commit sha which triggered the build
```
In ko builder this could be done with:
```
ko build --sbom=none --platform=linux/amd64 --tags=$VERSION --image-label="repo_url=${{ github.event.repository.html_url }},commit_sha=${{ github.sha }}"
```
See [ko-build.yaml](.github/workflows/ko-build.yaml) for more details

For docker build docker this could be done throug [dockerfile](https://dockerlabs.collabnix.com/beginners/dockerfile/Label_instruction.html) or [buildx](https://docs.docker.com/reference/cli/docker/buildx/build/)

## Configure devlake

Devlake allows different configurations of projects see [tutorial](https://devlake.apache.org/docs/Configuration/Tutorial)

The minimum is required:
- [Data connection](https://devlake.apache.org/docs/Overview/KeyConcepts/#data-connection)
- Data scope (under the details in data connection) - points what Git repose to track
- Scope config - describes what to fetch 
- Project - in a simpliest form on per data connection
- Webhook - endpoint to push deployments

## Setup handler

Modify `values.yaml` file (see schema for descriptions)
- Point to prometheus instance to read the data
- Define in `crane` section docker registry credentials and url to fetch image labels  
- `devlake` section to assign patterns of images to devlake projects.

Deploy the handler with updated helm values

## Test

Go to devlake grafana and check metrics - check `DORA Validation`, `DORA Details - Deployment Frequency` and `DORA Details - Lead Time for Changes` dashboards

If [handler-badges](https://github.com/tiktai/dora-badge) are installed check url
```
<handler-badge host>/v1/Tiktai-handler/df - to check are deployments pushed
<handler-badge host>/v1/Tiktai-handler/ltfc-stats - to check lead time for changes
```

## Running Locally
To start the server:
```bash
go run . --config ./.image-handler.yaml 
```

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

## License
See [LICENSE](LICENSE) for details.

---
*Generated on 2025-05-17*
