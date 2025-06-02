# Go Template Project

A ready-to-use Go template with DI, logging, configuration, database, HTTP server, observability, testing, and Docker support.

## Overview

This template provides a solid foundation for building Go applications with:

- **Dependency Injection**: Clean, testable architecture with proper DI container
- **Structured Logging**: Both structured and unstructured logging with Zap
- **Configuration Management**: Environment-based configuration with sensible defaults
- **Database Integration**: PostgreSQL and Kafka support with interfaces
- **HTTP Server**: Gin-based REST API with health checks and Swagger docs
- **Observability**: OpenTelemetry tracing and Prometheus metrics
- **Testing**: Comprehensive test coverage with proper mocking
- **Docker Support**: Production-ready containerization

## Getting Started

### Prerequisites

- Go 1.21+
- Docker (optional, for databases and monitoring)
- Make (for build automation)

### Configuration

The application uses environment variables with defaults:

```bash
# Application
GO_APP_APP_NAME=myapp
GO_APP_ENV=development
GO_APP_PORT=8080

# Logging
LOGGER_WRITE_STDOUT=true
LOGGER_ENABLE_STACK_TRACE=false
LOGGER_APP_LOG_PATH=./logs/app.log
LOGGER_ERR_LOG_PATH=./logs/error.log
```

## Architecture

### Project Structure

```
goapp/
├── api/                    # HTTP layer
│   ├── handlers/          # HTTP handlers with DI
│   └── routes/            # Route definitions
├── cmd/                   # Application entrypoints
│   ├── goapp/            # Main HTTP server
├── internal/              # Internal packages (not importable)
│   ├── config/           # Configuration management
│   ├── container/        # Dependency injection container
│   ├── db/              # Database implementations
│   │   ├── postgres/    # PostgreSQL client
│   │   └── kafka/       # Kafka client
│   ├── logging/         # Logging implementation
│   └── observability/   # OpenTelemetry setup
├── pkg/                  # Public packages (reusable)
│   └── clients/         # HTTP clients
└── docs/                # Swagger documentation
```

### Guide

- **pkg/ vs internal/**: `pkg/` contains reusable libraries, `internal/` contains app-specific code
- **Interface-driven**: All major components implement interfaces for easy testing
- **Dependency Injection**: No global state, all dependencies are injected
- **Configuration**: Single source of truth with environment variable support
- **Error Handling**: Proper error propagation with context

## Features

### Logging

The logging package supports both structured and unstructured logging:

```go
// Structured logging (recommended)
logger.Info("User created", 
    logging.String("user_id", "123"),
    logging.String("email", "user@example.com"))

// Unstructured logging (for simple cases)
logger.Infof("User %s created with email %s", userID, email)

// Context logging
userLogger := logger.With(logging.String("user_id", userID))
userLogger.Info("Processing user request")
```

### Database

PostgreSQL integration with proper error handling:

```go
// The database interface is injected into handlers
func (h *Handler) CreateUser(c *gin.Context) {
    // Use h.Database.DB() to get *sqlx.DB instance
    user := &User{}
    err := h.Database.DB().Get(user, "SELECT * FROM users WHERE id = $1", userID)
    if err != nil {
        h.Logger.Error("Failed to fetch user", logging.Error(err))
        return
    }
}
```

### Health Checks

Built-in health checks that verify:

- Application status
- Database connectivity
- External service availability

### Configuration

Type-safe configuration with validation:

```go
type AppConfig struct {
    Name        string `envconfig:"APP_NAME" default:"goapp"`
    Env         string `envconfig:"ENV" default:"development"`
    Port        int    `envconfig:"PORT" default:"8080"`
}
```

## Development

### Building

```bash
# Build for current platform
go build ./cmd/goapp

# Build for Linux
GOOS=linux go build ./cmd/goapp

# Build with race detector
go build -race ./cmd/goapp
```

### Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with race detector
go test -race ./...
```

### Docker

```bash
# Build image
docker build -t myapp .

# Run with docker-compose
docker-compose up
```

## Extending the Template

### Adding New Dependencies

1. Update `internal/config/config.go` with new configuration
2. Add the service interface to your package
3. Update `internal/container/container.go` to initialize the new service
4. Inject into handlers via the container

### Adding New Endpoints

1. Create handler methods in `api/handlers/`
2. Update `api/routes/routes.go` to register routes
3. Add Swagger documentation comments
4. Generate docs: `swag init`

### Adding Middleware

```go
// In routes.go
router.Use(authMiddleware(container))
router.Use(corsMiddleware())
```

## Production Considerations

- Set `GO_APP_ENV=production` for optimized builds
- Use proper secret management for sensitive configuration
- Set up log aggregation (ELK stack, Fluentd, etc.)
- Configure reverse proxy (nginx, Traefik)
- Set up monitoring and alerting
- Use database migrations for schema changes
- Configure graceful shutdown timeouts appropriately

## Monitoring

The template includes:

- **Prometheus metrics**: Available at `/metrics`
- **OpenTelemetry tracing**: Distributed tracing support
- **Health checks**: Kubernetes-ready health endpoints
- **Structured logging**: JSON formatted logs for aggregation

## License

This template is provided as-is for creating new Go projects. Modify as needed for your use case.
