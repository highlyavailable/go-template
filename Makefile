APP_NAME=goapp
CONTAINER_NAME=$(APP_NAME)
IMAGE_NAME=$(APP_NAME)
CD_APP=cd $(APP_NAME) &&

BINARY_PATH=./build/$(APP_NAME)
ENTRY_POINT=./cmd/$(APP_NAME)
PATH_TO_DOCKER_RECYCLE=./docker-recycle.sh

# Set the path to the .env file, defaulting to the current directory if not set
ENV_PATH ?= .env

# Load environment variables
ifneq (,$(wildcard $(ENV_PATH)))
    include $(ENV_PATH)
    export
endif

all: env swagger build run

env:
	@echo "-> Checking for .env file"
	@if [ ! -f $(ENV_PATH) ]; then \
        echo ".env file not found, creating one"; \
        touch $(ENV_PATH); \
        echo "APP_NAME=$(APP_NAME)" >> $(ENV_PATH); \
        echo "ENV=development" >> $(ENV_PATH); \
        echo "PORT=8080" >> $(ENV_PATH); \
        echo "LOG_LEVEL=debug" >> $(ENV_PATH); \
        echo ".env file created"; \
	else \
		echo ".env file found"; \
	fi
	set -a; source $(ENV_PATH); set +a

swagger:
	@echo "-> Generating swagger docs"
	@echo "Working directory: $(CD_APP)"
	$(CD_APP) swag init -g $(ENTRY_POINT)/main.go

build:
	@echo "-> Building $(APP_NAME)"
	$(CD_APP) go build -o $(BINARY_PATH) $(ENTRY_POINT)

build-verbose:
	@echo "-> Building $(APP_NAME)"
	$(CD_APP) go build -v -o $(BINARY_PATH) $(ENTRY_POINT)

build-linux:
	@echo "-> Building $(APP_NAME) for linux"
	$(CD_APP) GOOS=linux GOARCH=amd64 go build -o $(BINARY_PATH) $(ENTRY_POINT)

build-windows:
	@echo "-> Building $(APP_NAME) for windows"
	$(CD_APP) GOOS=windows GOARCH=amd64 go build -o $(BINARY_PATH).exe $(ENTRY_POINT)

build-race:
	@echo "-> Building $(APP_NAME) with race detector"
	$(CD_APP) go build -race -o $(BINARY_PATH) $(ENTRY_POINT)

run:
	@echo "-> Running $(APP_NAME)"
	$(CD_APP) chmod +x $(BINARY_PATH)
	$(CD_APP) $(BINARY_PATH)

docker-recycle:
	@echo "-> Recycle docker containers"
	chmod +x $(PATH_TO_DOCKER_RECYCLE)
	$(PATH_TO_DOCKER_RECYCLE)

docker-exec:
	@echo "-> Executing shell in $(CONTAINER_NAME)"
	docker run -it --entrypoint /bin/sh -e ENV_PATH=$(CONTAINER_ENV_PATH) $(CONTAINER_NAME)

tidy:
	@echo "-> Running go mod tidy"
	$(CD_APP) go mod tidy