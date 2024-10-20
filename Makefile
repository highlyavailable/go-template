APP_NAME := goapp
CONTAINER_NAME := $(APP_NAME)
IMAGE_NAME := $(APP_NAME)

BINARY_PATH := $(APP_NAME)/build/$(APP_NAME)
ENTRY_POINT := $(APP_NAME)/cmd/$(APP_NAME)
PATH_TO_DOCKER_RECYCLE := ./docker-recycle.sh

ENV_FILE ?= .env

# Load environment variables
ifneq (,$(wildcard $(ENV_FILE)))
    include $(ENV_FILE)
    export $(shell sed 's/=.*//' $(ENV_FILE))
endif

# Define colors
CYAN := \033[36m
RESET := \033[0m

# Define a function for prettier logging
define log
	@echo "$(CYAN)$(1)$(RESET)"
endef

env:
	$(call log,Checking for .env file)
	@if [ ! -f $(ENV_FILE) ]; then \
		echo "$(CYAN).env file not found, creating one$(RESET)"; \
		touch $(ENV_FILE); \
	else \
		echo "$(CYAN).env file found$(RESET)"; \
	fi

tidy:
	$(call log,Running go mod tidy)
	cd $(APP_NAME) && go mod tidy

lint:
	$(call log,Running golint)
	cd $(APP_NAME) && golint ./...

test:
	$(call log,Running tests)
	cd $(APP_NAME) && go test -v ./...

test-coverage:
	$(call log,Running tests with coverage)
	cd $(APP_NAME) && go test -v -coverprofile=coverage.out ./...
	cd $(APP_NAME) && go tool cover -html=coverage.out -o coverage.html

security-scan:
	$(call log,Running security scan with gosec)
	cd $(APP_NAME) && gosec ./...

swagger:
	$(call log,Generating swagger docs)
	cd $(APP_NAME) && swag init -g cmd/$(APP_NAME)/main.go

build:
	$(call log,Building $(APP_NAME))
	cd $(APP_NAME) && go build -o build/$(APP_NAME) cmd/$(APP_NAME)/main.go

build-verbose:
	$(call log,Building $(APP_NAME) with verbose output)
	cd $(APP_NAME) && go build -v -o build/$(APP_NAME) cmd/$(APP_NAME)/main.go

run:
	$(call log,Running $(APP_NAME))
	chmod +x $(BINARY_PATH)
	$(BINARY_PATH)

monitors:
	$(call log,Starting monitoring stack)
	docker-compose -f monitoring/docker-compose.yml up -d

prevalidate: env lint tidy
all: prevalidate build run
ci: prevalidate test-coverage build-verbose run
gin: prevalidate swagger build run

.PHONY: env lint test test-coverage security-scan swagger build build-verbose run prevalidate all ci