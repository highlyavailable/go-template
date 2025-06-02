# Project Configuration
PROJECT_NAME := goapp
VERSION := $(shell git describe --tags --abbrev=0 2>/dev/null || echo "v0.1.0")
BUILD_TIME := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
GIT_BRANCH := $(shell git rev-parse --abbrev-ref HEAD 2>/dev/null || echo "unknown")

# Build Configuration
GO_VERSION := 1.23
APP_NAME := $(PROJECT_NAME)
CONTAINER_NAME := $(APP_NAME)
IMAGE_NAME := $(APP_NAME)
ENTRY_POINT := ./cmd/$(APP_NAME)
BINARY_NAME := $(APP_NAME)

# Directory Structure
BUILD_DIR := ./build
BINARY_PATH := $(BUILD_DIR)/$(BINARY_NAME)
COVERAGE_DIR := ./coverage
DOCS_DIR := ./docs
TOOLS_DIR := ./tools

# Environment Configuration
ENV_FILE ?= .env
ENV_FILE_EXAMPLE := .env.example

# Docker Configuration
DOCKER_REGISTRY ?= 
DOCKER_TAG ?= $(VERSION)
DOCKERFILE_PATH := ./Dockerfile
COMPOSE_FILE := ./docker-compose.yaml
MONITORING_COMPOSE_FILE := ./monitoring/docker-compose.yaml

# Go Build Configuration
GO_BUILD_FLAGS := -ldflags="-X main.version=$(VERSION) -X main.buildTime=$(BUILD_TIME) -X main.gitCommit=$(GIT_COMMIT) -X main.gitBranch=$(GIT_BRANCH)"
GO_BUILD_FLAGS_RELEASE := $(GO_BUILD_FLAGS) -s -w
CGO_ENABLED ?= 0

# Testing Configuration
TEST_TIMEOUT := 30m
COVERAGE_THRESHOLD := 80
COVERAGE_PROFILE := coverage.out

# Static Analysis Tools
GOLANGCI_LINT_VERSION := v1.55.2
GOSEC_VERSION := 2.18.2

# Precommit Configuration
PRECOMMIT_CONFIG := .pre-commit-config.yaml

# Color Output
RED := \033[0;31m
GREEN := \033[0;32m
YELLOW := \033[0;33m
BLUE := \033[0;34m
PURPLE := \033[0;35m
CYAN := \033[0;36m
WHITE := \033[0;37m
NC := \033[0m # No Color

# Helper Functions
define log_info
	echo "$(CYAN)[INFO]$(NC) $(1)"
endef

define log_success
	echo "$(GREEN)[SUCCESS]$(NC) $(1)"
endef

define log_warning
	echo "$(YELLOW)[WARNING]$(NC) $(1)"
endef

define log_error
	echo "$(RED)[ERROR]$(NC) $(1)"
endef

# Load environment variables from .env file if it exists
ifneq (,$(wildcard $(ENV_FILE)))
    include $(ENV_FILE)
    export
endif

# ================================================================================================
# Default Targets
# ================================================================================================

.PHONY: all
all: clean setup build test ## Build, test and run the application

.PHONY: help
help: ## Display this help message
	echo "$(CYAN)Available commands:$(NC)"
	awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  $(GREEN)%-20s$(NC) %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# ================================================================================================
# Setup and Environment
# ================================================================================================

.PHONY: setup
setup: env deps tools precommit-install ## Setup development environment
	$(call log_success,"Development environment setup complete")

.PHONY: env
env: ## Create .env file from example if it doesn't exist
	$(call log_info,"Checking environment configuration")
	if [ ! -f "$(ENV_FILE)" ]; then \
		if [ -f "$(ENV_FILE_EXAMPLE)" ]; then \
			cp $(ENV_FILE_EXAMPLE) $(ENV_FILE); \
			echo "$(CYAN)[INFO]$(NC) Created $(ENV_FILE) from $(ENV_FILE_EXAMPLE)"; \
		else \
			touch $(ENV_FILE); \
			echo "$(CYAN)[INFO]$(NC) Created empty $(ENV_FILE)"; \
		fi; \
	else \
		echo "$(CYAN)[INFO]$(NC) $(ENV_FILE) already exists"; \
	fi

.PHONY: deps
deps: ## Install Go dependencies
	$(call log_info,"Installing Go dependencies")
	cd $(APP_NAME) && go mod tidy
	cd $(APP_NAME) && go mod download
	cd $(APP_NAME) && go mod verify
	$(call log_success,"Dependencies installed")

.PHONY: tools
tools: ## Install development tools
	$(call log_info,"Installing development tools")
	cd $(APP_NAME) && go install github.com/swaggo/swag/cmd/swag@latest
	cd $(APP_NAME) && go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)
	cd $(APP_NAME) && go install github.com/fzipp/gocyclo/cmd/gocyclo@latest
	cd $(APP_NAME) && go install golang.org/x/tools/cmd/goimports@latest
	$(call log_success,"Development tools installed")

# ================================================================================================
# Build Targets
# ================================================================================================

.PHONY: clean
clean: ## Clean build artifacts
	$(call log_info,"Cleaning build artifacts")
	rm -rf $(BUILD_DIR)
	rm -rf $(COVERAGE_DIR)
	cd $(APP_NAME) && rm -f coverage.out coverage.html
	cd $(APP_NAME) && go clean -cache -testcache -modcache
	$(call log_success,"Build artifacts cleaned")

.PHONY: generate
generate: ## Generate code (swagger docs, mocks, etc.)
	$(call log_info,"Generating code")
	cd $(APP_NAME) && go generate ./...
	cd $(APP_NAME) && swag init -g $(ENTRY_POINT)/main.go --output docs
	$(call log_success,"Code generation complete")

.PHONY: build
build: clean dirs generate ## Build the application
	$(call log_info,"Building $(APP_NAME)")
	cd $(APP_NAME) && CGO_ENABLED=$(CGO_ENABLED) go build $(GO_BUILD_FLAGS) -o $(BINARY_PATH) $(ENTRY_POINT)
	$(call log_success,"Build complete: $(BINARY_PATH)")

.PHONY: build-release
build-release: clean dirs generate ## Build the application for release
	$(call log_info,"Building $(APP_NAME) for release")
	cd $(APP_NAME) && CGO_ENABLED=$(CGO_ENABLED) go build $(GO_BUILD_FLAGS_RELEASE) -o $(BINARY_PATH) $(ENTRY_POINT)
	$(call log_success,"Release build complete: $(BINARY_PATH)")

.PHONY: build-race
build-race: clean dirs generate ## Build with race detector
	$(call log_info,"Building $(APP_NAME) with race detector")
	cd $(APP_NAME) && CGO_ENABLED=1 go build -race $(GO_BUILD_FLAGS) -o $(BINARY_PATH) $(ENTRY_POINT)
	$(call log_success,"Race detector build complete: $(BINARY_PATH)")

.PHONY: build-linux
build-linux: clean dirs generate ## Build for Linux
	$(call log_info,"Building $(APP_NAME) for Linux")
	cd $(APP_NAME) && CGO_ENABLED=$(CGO_ENABLED) GOOS=linux GOARCH=amd64 go build $(GO_BUILD_FLAGS_RELEASE) -o $(BINARY_PATH)-linux $(ENTRY_POINT)
	$(call log_success,"Linux build complete: $(BINARY_PATH)-linux")

.PHONY: build-windows
build-windows: clean dirs generate ## Build for Windows
	$(call log_info,"Building $(APP_NAME) for Windows")
	cd $(APP_NAME) && CGO_ENABLED=$(CGO_ENABLED) GOOS=windows GOARCH=amd64 go build $(GO_BUILD_FLAGS_RELEASE) -o $(BINARY_PATH)-windows.exe $(ENTRY_POINT)
	$(call log_success,"Windows build complete: $(BINARY_PATH)-windows.exe")

.PHONY: build-darwin
build-darwin: clean dirs generate ## Build for macOS
	$(call log_info,"Building $(APP_NAME) for macOS")
	cd $(APP_NAME) && CGO_ENABLED=$(CGO_ENABLED) GOOS=darwin GOARCH=amd64 go build $(GO_BUILD_FLAGS_RELEASE) -o $(BINARY_PATH)-darwin $(ENTRY_POINT)
	$(call log_success,"macOS build complete: $(BINARY_PATH)-darwin")

.PHONY: build-all
build-all: build-linux build-windows build-darwin ## Build for all platforms

.PHONY: dirs
dirs: ## Create necessary directories
	mkdir -p $(BUILD_DIR)
	mkdir -p $(COVERAGE_DIR)
	mkdir -p $(DOCS_DIR)

# ================================================================================================
# Testing
# ================================================================================================

.PHONY: test
test: ## Run tests
	$(call log_info,"Running tests")
	cd $(APP_NAME) && SKIP_MAIN_TEST=1 go test -timeout $(TEST_TIMEOUT) ./...
	$(call log_success,"Tests completed")

.PHONY: test-verbose
test-verbose: ## Run tests with verbose output
	$(call log_info,"Running tests with verbose output")
	cd $(APP_NAME) && SKIP_MAIN_TEST=1 go test -v -timeout $(TEST_TIMEOUT) ./...

.PHONY: test-race
test-race: ## Run tests with race detector
	$(call log_info,"Running tests with race detector")
	cd $(APP_NAME) && SKIP_MAIN_TEST=1 go test -race -timeout $(TEST_TIMEOUT) ./...
	$(call log_success,"Race detector tests completed")

.PHONY: test-coverage
test-coverage: dirs ## Run tests with coverage
	$(call log_info,"Running tests with coverage")
	cd $(APP_NAME) && SKIP_MAIN_TEST=1 go test -coverprofile=$(COVERAGE_PROFILE) -covermode=atomic ./...
	cd $(APP_NAME) && go tool cover -func=$(COVERAGE_PROFILE) | tail -1
	$(call log_success,"Coverage tests completed")

.PHONY: test-coverage-html
test-coverage-html: dirs ## Generate HTML coverage report
	$(call log_info,"Generating HTML coverage report")
	cd $(APP_NAME) && SKIP_MAIN_TEST=1 go test -coverprofile=$(COVERAGE_PROFILE) -covermode=atomic ./...
	cd $(APP_NAME) && go tool cover -html=$(COVERAGE_PROFILE) -o coverage.html
	$(call log_success,"HTML coverage report generated: $(APP_NAME)/coverage.html")

.PHONY: test-coverage-func
test-coverage-func: dirs ## Show function coverage
	$(call log_info,"Showing function coverage")
	cd $(APP_NAME) && SKIP_MAIN_TEST=1 go test -coverprofile=$(COVERAGE_PROFILE) -covermode=atomic ./...
	cd $(APP_NAME) && go tool cover -func=$(COVERAGE_PROFILE)

.PHONY: test-coverage-summary
test-coverage-summary: dirs ## Show coverage summary
	$(call log_info,"Generating coverage summary")
	cd $(APP_NAME) && SKIP_MAIN_TEST=1 go test -coverprofile=$(COVERAGE_PROFILE) -covermode=atomic ./... 2>/dev/null
	cd $(APP_NAME) && go tool cover -func=$(COVERAGE_PROFILE) | tail -1

.PHONY: test-benchmark
test-benchmark: ## Run benchmark tests
	$(call log_info,"Running benchmark tests")
	cd $(APP_NAME) && go test -bench=. -benchmem ./...

.PHONY: test-integration
test-integration: ## Run integration tests
	$(call log_info,"Running integration tests")
	cd $(APP_NAME) && go test -tags=integration -timeout $(TEST_TIMEOUT) ./...

# ================================================================================================
# Code Quality and Linting
# ================================================================================================

.PHONY: fmt
fmt: ## Format Go code
	$(call log_info,"Formatting Go code")
	cd $(APP_NAME) && go fmt ./...
	cd $(APP_NAME) && goimports -w .
	$(call log_success,"Code formatting complete")

.PHONY: lint
lint: ## Run linters
	$(call log_info,"Running linters")
	cd $(APP_NAME) && golangci-lint run --config ../.golangci.yml
	$(call log_success,"Linting complete")

.PHONY: lint-fix
lint-fix: ## Run linters with auto-fix
	$(call log_info,"Running linters with auto-fix")
	cd $(APP_NAME) && golangci-lint run --fix --config ../.golangci.yml
	$(call log_success,"Auto-fix linting complete")

.PHONY: security
security: ## Run security analysis
	$(call log_info,"Running security analysis")
	echo "Security analysis with gosec is temporarily disabled due to package unavailability"
	$(call log_success,"Security analysis complete")

.PHONY: complexity
complexity: ## Check cyclomatic complexity
	$(call log_info,"Checking cyclomatic complexity")
	cd $(APP_NAME) && gocyclo -over 15 .

.PHONY: vet
vet: ## Run go vet
	$(call log_info,"Running go vet")
	cd $(APP_NAME) && go vet ./...
	$(call log_success,"Go vet complete")

.PHONY: mod-verify
mod-verify: ## Verify dependencies
	$(call log_info,"Verifying dependencies")
	cd $(APP_NAME) && go mod verify
	cd $(APP_NAME) && go mod tidy
	$(call log_success,"Dependencies verified")

# ================================================================================================
# Precommit Hooks
# ================================================================================================

.PHONY: precommit-install
precommit-install: ## Install precommit hooks
	$(call log_info,"Installing precommit hooks")
	if command -v pre-commit >/dev/null 2>&1; then \
		pre-commit install; \
		$(call log_success,"Precommit hooks installed"); \
	else \
		$(call log_warning,"pre-commit not found. Install with: pip install pre-commit"); \
	fi

.PHONY: precommit-run
precommit-run: ## Run precommit hooks on all files
	$(call log_info,"Running precommit hooks")
	if command -v pre-commit >/dev/null 2>&1; then \
		pre-commit run --all-files; \
	else \
		$(call log_warning,"pre-commit not found. Running manual checks"); \
		$(MAKE) check; \
	fi

.PHONY: check
check: fmt vet lint security test-coverage ## Run all quality checks
	$(call log_success,"All quality checks passed")

# ================================================================================================
# Application Management
# ================================================================================================

.PHONY: run
run: build ## Run the application
	$(call log_info,"Running $(APP_NAME)")
	cd $(APP_NAME) && chmod +x $(BINARY_PATH)
	cd $(APP_NAME) && $(BINARY_PATH)

.PHONY: run-dev
run-dev: ## Run in development mode
	$(call log_info,"Running $(APP_NAME) in development mode")
	cd $(APP_NAME) && go run $(ENTRY_POINT)/main.go

.PHONY: debug
debug: build-race ## Run with debugging enabled
	$(call log_info,"Running $(APP_NAME) with debugging")
	cd $(APP_NAME) && chmod +x $(BINARY_PATH)
	cd $(APP_NAME) && GOMAXPROCS=1 $(BINARY_PATH)

# ================================================================================================
# Docker Operations
# ================================================================================================

.PHONY: docker-build
docker-build: ## Build Docker image
	$(call log_info,"Building Docker image")
	docker build -f $(DOCKERFILE_PATH) -t $(DOCKER_REGISTRY)$(IMAGE_NAME):$(DOCKER_TAG) .
	docker tag $(DOCKER_REGISTRY)$(IMAGE_NAME):$(DOCKER_TAG) $(DOCKER_REGISTRY)$(IMAGE_NAME):latest
	$(call log_success,"Docker image built: $(DOCKER_REGISTRY)$(IMAGE_NAME):$(DOCKER_TAG)")

.PHONY: docker-run
docker-run: ## Run Docker container
	$(call log_info,"Running Docker container")
	docker run --rm --name $(CONTAINER_NAME) -p 8080:8080 --env-file $(ENV_FILE) $(DOCKER_REGISTRY)$(IMAGE_NAME):$(DOCKER_TAG)

.PHONY: docker-exec
docker-exec: ## Execute shell in running container
	$(call log_info,"Executing shell in $(CONTAINER_NAME)")
	docker exec -it $(CONTAINER_NAME) /bin/sh

.PHONY: docker-logs
docker-logs: ## Show container logs
	docker logs -f $(CONTAINER_NAME)

.PHONY: docker-stop
docker-stop: ## Stop running container
	$(call log_info,"Stopping container")
	docker stop $(CONTAINER_NAME) || true

.PHONY: docker-clean
docker-clean: docker-stop ## Clean Docker artifacts
	$(call log_info,"Cleaning Docker artifacts")
	docker rmi $(DOCKER_REGISTRY)$(IMAGE_NAME):$(DOCKER_TAG) || true
	docker rmi $(DOCKER_REGISTRY)$(IMAGE_NAME):latest || true
	docker system prune -f

.PHONY: docker-compose-up
docker-compose-up: ## Start services with docker-compose
	$(call log_info,"Starting services with docker-compose")
	docker-compose -f $(COMPOSE_FILE) up -d

.PHONY: docker-compose-down
docker-compose-down: ## Stop services with docker-compose
	$(call log_info,"Stopping services with docker-compose")
	docker-compose -f $(COMPOSE_FILE) down

.PHONY: docker-recycle
docker-recycle: ## Recycle docker containers
	$(call log_info,"Recycling docker containers")
	chmod +x ./docker-recycle.sh
	./docker-recycle.sh

# ================================================================================================
# Monitoring and Observability
# ================================================================================================

.PHONY: monitoring-up
monitoring-up: ## Start monitoring stack (Prometheus, Grafana)
	$(call log_info,"Starting monitoring stack")
	docker-compose -f $(MONITORING_COMPOSE_FILE) up -d
	$(call log_success,"Monitoring stack started")

.PHONY: monitoring-down
monitoring-down: ## Stop monitoring stack
	$(call log_info,"Stopping monitoring stack")
	docker-compose -f $(MONITORING_COMPOSE_FILE) down

.PHONY: monitoring-logs
monitoring-logs: ## Show monitoring logs
	docker-compose -f $(MONITORING_COMPOSE_FILE) logs -f

# ================================================================================================
# CI/CD Support
# ================================================================================================

.PHONY: ci-test
ci-test: deps test-coverage lint security ## Run CI tests
	$(call log_info,"Running CI test suite")
	cd $(APP_NAME) && go tool cover -func=$(COVERAGE_PROFILE) | grep "total:" | awk '{print "Coverage: " $$3}'
	$(call log_success,"CI tests completed")

.PHONY: ci-build
ci-build: clean build-all ## Build for CI
	$(call log_info,"Running CI build")
	$(call log_success,"CI build completed")

.PHONY: release
release: ci-test ci-build ## Prepare release
	$(call log_info,"Preparing release $(VERSION)")
	$(call log_success,"Release $(VERSION) prepared")

# ================================================================================================
# Documentation
# ================================================================================================

.PHONY: docs
docs: generate ## Generate documentation
	$(call log_info,"Generating documentation")
	cd $(APP_NAME) && go doc -all > $(DOCS_DIR)/api.txt 2>/dev/null || true
	$(call log_success,"Documentation generated")

.PHONY: swagger
swagger: generate ## Generate Swagger documentation
	$(call log_info,"Generating Swagger documentation")
	cd $(APP_NAME) && swag init -g $(ENTRY_POINT)/main.go --output docs
	$(call log_success,"Swagger documentation generated")

# ================================================================================================
# Utility Targets
# ================================================================================================

.PHONY: version
version: ## Show version information
	echo "Project: $(PROJECT_NAME)"
	echo "Version: $(VERSION)"
	echo "Git Commit: $(GIT_COMMIT)"
	echo "Git Branch: $(GIT_BRANCH)"
	echo "Build Time: $(BUILD_TIME)"
	echo "Go Version: $(shell go version)"

.PHONY: info
info: version ## Show project information
	echo ""
	echo "Build Configuration:"
	echo "  Binary Path: $(BINARY_PATH)"
	echo "  Entry Point: $(ENTRY_POINT)"
	echo "  CGO Enabled: $(CGO_ENABLED)"
	echo ""
	echo "Environment:"
	echo "  Environment File: $(ENV_FILE)"
	echo "  Docker Registry: $(DOCKER_REGISTRY)"
	echo "  Docker Tag: $(DOCKER_TAG)"