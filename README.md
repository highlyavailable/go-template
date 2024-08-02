# Go Template Project

## Overview

Fork for quick creation of new Go applications. It includes a structured directory layout with a Makefile for build automation, `envconfig` for environment configuration/management, `Zap` and `lumberjack` for structured logging with rotation and retention, `OpenTelemetry` for observability and tracing, `Prometheus` for metrics collection and monitoring. Includes a Dockerfile and docker-compose.yaml for containerization. Contains a basic HTTP Gin server with health check and Swagger documentation, with additional integrations for Kafka, Postgres, MongoDB, and Redis.

## Getting Started

### Prerequisites

- Go 1.21 (specified in go.mod)
- make installed
- Docker (optional)

### Setup

1. **Clone the Repository**

   ```bash
   git clone <repository-url>
   cd <repository-name>
   ```

2. **Install Dependencies**

   Install the project dependencies using `go mod`:

   ```bash
   cd goapp
   go mod tidy
   ```

3. **Rename the Application**

   Update the application name using the `rename_app.sh` script.

   ```bash
   ./rename_app.sh <old_app_name> <new_app_name>
   # Example
   ./rename_app.sh goapp myapp
   ```

4. **Set Up Environment Variables**

   Ensure that required environment variables are set. Refer to the `.env.example` file for reference.


### Makefile Targets

The Makefile in the root of the project provides a set of commands to help with building, running, and managing the application. Running `make first` will set up your project by creating necessary files and installing dependencies.

#### Setup

- **`first`**: Performs an initial setup by:
  - Creating the `.env` file if it doesn’t exist.
  - Installing Swaggo for Swagger documentation.
  - Generating Swagger documentation.
  - Building the application.
  - Running the application.

#### Environment Management

- **`env`**: Checks for the `.env` file and creates it if missing.

#### Swagger Documentation

- **`swaggo`**: Installs the Swaggo tool for generating Swagger documentation.
- **`swagger`**: Generates Swagger documentation from the codebase.

#### Build

- **`build`**: Builds the application binary defaults to MacOS.
- **`build-verbose`**: Builds the application with verbose output.
- **`build-linux`**: Builds the application for Linux.
- **`build-windows`**: Builds the application for Windows.
- **`build-race`**: Builds the application with race detector enabled.

#### Run

- **`run`**: Runs the built application binary, stored in the `build` directory.

#### Docker Management

- **`docker-recycle`**: Recycles Docker containers using the `docker-recycle.sh` script, located in the root of the project.
- **`docker-exec`**: Executes a shell in the Docker container.

#### Miscellaneous

- **`tidy`**: Tidy will clean up the Go module cache in the Go module directory.


Running `make first` is recommended to get your project set up and ready for development.

### Docker

To run the application in Docker:

```bash
docker-compose up --build
```

Attach a shell to the container:

```bash
docker run -it --entrypoint /bin/sh go-app
```

## File Structure

```bash
.
├── Dockerfile
├── Makefile
├── README.md
├── assets
│   └── certs
├── docker-compose.yaml
├── docker-recycle.sh
├── goapp
│   ├── api
│   │   ├── handlers
│   │   │   └── health.go
│   │   └── routes
│   │       └── routes.go
│   ├── build
│   │   └── goapp
│   ├── cmd
│   │   └── goapp
│   │       └── main.go
│   ├── docs
│   │   ├── docs.go
│   │   ├── swagger.json
│   │   └── swagger.yaml
│   ├── go.mod
│   ├── go.sum
│   ├── internal
│   │   └── config
│   │       └── config.go
│   ├── pkg
│   │   ├── clients
│   │   │   └── clients.go
│   │   ├── db
│   │   │   ├── kafka
│   │   │   │   ├── example.go
│   │   │   │   ├── kafka.go
│   │   │   │   └── model.go
│   │   │   └── postgres
│   │   │       ├── example.go
│   │   │       ├── model.go
│   │   │       └── postgres.go
│   │   ├── logging
│   │   │   ├── logging.go
│   │   │   └── model.go
│   │   └── otel
│   │       └── otel.go
│   ├── scripts
│   ├── tests
│   │   ├── helpers
│   │   └── unit
│   └── web
│       ├── static
│       └── templates
├── logs
│   └── app.log
├── monitoring
│   ├── datasources
│   │   └── datasources.yml
│   ├── docker-compose.yaml
│   └── prometheus.yaml
└── rename_app.sh
```

## Directory Overview

- **api**: Contains the HTTP handlers and routes.
- **build**: Build artifacts.
- **cmd**: Main applications for the project.
- **docs**: Documentation including Swagger.
- **internal**: Internal application code.
- **pkg**: Library code to be used by external applications.
  - **db/kafka**: Kafka integration.
  - **db/postgres**: PostgreSQL integration.
  - **logging**: Logging setup.
  - **otel**: OpenTelemetry setup.
- **scripts**: Scripts for various tasks.
- **tests**: Unit and helper tests.
- **web**: Web assets like static files and templates.
- **logs**: Application logs.
- **monitoring**: Monitoring setup for Grafana and Prometheus.

## Environment Variables

The application uses `envconfig` to manage environment variables. Ensure the following variables are set:

- **Kafka**:
  - `KAFKA_BROKERS`
  - `KAFKA_PRODUCER_TOPIC`
  - `KAFKA_CONSUMER_TOPIC`
  - `KAFKA_CONSUMER_GROUP`
  - `KAFKA_CONSUMER_OFFSET`

- **Postgres**:
  - `DB_HOST`
  - `DB_PORT`
  - `DB_USER`
  - `DB_PASSWORD`
  - `DB_NAME`
  - `DB_SSL_MODE`

- **Prometheus**:
  - `PROMETHEUS_PORT`

## Logging

The application uses Zap for structured logging with log rotation and retention. Logs are written to `app.log` and `error.log` in the `logs` directory.

## Monitoring

Prometheus is used for metrics collection and Grafana for visualization. The monitoring setup is available in the `monitoring` directory.