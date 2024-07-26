# Go Template App

## Overview

Go Template App is a Go application template with a structured layout and Docker support. Use the `Makefile` to build and run the application efficiently.

## Getting Started

### Prerequisites

Ensure you have `make` and `Go` installed. Docker is optional for containerized runs.

### Build and Run

1. **Rename the Application**

   Update the application name using the `rename_app.sh` script. If cloning the repository, the default application name is `goapp`.

   ```bash
   ./rename_app.sh <old_app_name> <new_app_name> # Default old_app_name is goapp
   
   # Example
    ./rename_app.sh goapp myapp
   ```

2. **Build and Run**

   Use `make` to build and run the application:

   ```bash
   make
   ```

   - **Build Variants**:
     - `make build` – Build the application.
     - `make build-verbose` – Build with verbose output.
     - `make build-linux` – Build for Linux.
     - `make build-windows` – Build for Windows.
     - `make build-race` – Build with race detector enabled.

   - **Run**:
     - The `make` command automatically builds and runs the application.
   
3. **Attach a Shell to the Container**

   If using Docker, run a shell in the container:

   ```bash
   docker run -it --entrypoint /bin/sh go-template-app
   ```

## File Structure

```
.
├── Dockerfile
├── Makefile
├── README.md
├── app
│   ├── build
│   │   └── app
│   └── cmd
├── docker-compose.yml
├── goapp
│   ├── build
│   │   └── goapp
│   ├── cmd
│   │   └── goapp
│   │       └── main.go
│   ├── configs
│   │   └── config.yaml
│   ├── docs
│   ├── go.mod
│   ├── go.sum
│   ├── internal
│   │   ├── config
│   │   │   └── config.go
│   │   └── models
│   ├── pkg
│   └── scripts
├── rename_app.sh
```