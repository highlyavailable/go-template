# Go Template Project

A robust Go application template with integrated observability, persistence, and message queue support.

## Quick Start

```bash
make first  # Initial setup
make run    # Run the application
````

## Key Features

- Structured logging (Zap + lumberjack)
- Observability (OpenTelemetry)
- Metrics (Prometheus)
- HTTP server (Gin) with Swagger docs
- Integrations: Kafka, Postgres, MongoDB, Redis
- Docker and Docker Compose support

## Development

### Prerequisites

- Go 1.21+
- Docker and Docker Compose

### Make Targets

- `make build`: Build the application (output in `goapp/build`)
- `make run`: Run the application
- `make test`: Run tests
- `make test-coverage`: Run tests with coverage (generates `goapp/coverage.html`)
- `make lint`: Run golint
- `make security-scan`: Run gosec for security analysis
- `make swagger`: Generate Swagger documentation
- `make monitors`: Start monitoring stack (Prometheus + Grafana)

### CI Pipeline

````bash
make ci  # Runs linting, tests with coverage, security scan, and verbose build
````

## Configuration

Environment variables are managed via `envconfig`. See `.env.example` for required variables.

## Docker

The project includes a `Dockerfile` and `docker-compose.yaml` for containerized development and deployment.

````bash
docker compose up -d  # Start all services
````

## Monitoring

Prometheus and Grafana configurations are in the `monitoring` directory.

````bash
make monitors  # Starts the monitoring stack
````

## Directory Structure

````
.
├── goapp/
│   ├── api/          # HTTP handlers and routes
│   ├── build/        # Compiled binaries
│   ├── cmd/          # Main applications
│   ├── internal/     # Private application code
│   ├── pkg/          # Public libraries
│   └── web/          # Web assets
├── monitoring/       # Monitoring configurations
├── Dockerfile
├── docker-compose.yaml
├── Makefile
└── README.md
````

For detailed information on each component, refer to the respective package documentation.
