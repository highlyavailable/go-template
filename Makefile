APP_NAME=goapp
BINARY_PATH=./build/$(APP_NAME)
ENTRY_POINT=./cmd/$(APP_NAME)
CD_APP=cd $(APP_NAME) &&

all: env build run

env:
	@echo "-> Checking for .env file"
	@if [ ! -f .env ]; then \
	    echo ".env file not found, creating one"; \
	    cp .env.example .env; \
	fi
	set -a; source .env; set +a

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