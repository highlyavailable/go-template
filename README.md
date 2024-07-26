# Go Template

## Overview

A Go application template with a structured layout and Docker support. `Makefile` is included for easy build and run commands.

## Getting Started

### Prerequisites

- Go 1.20 (specified in `go.mod`)
- `make` installed
- Docker (optional, for containerized runs)

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
Tree structure of the template:

```
> tree
.
├── Dockerfile
├── Makefile
├── README.md
├── docker-compose.yml
├── goapp
│   ├── assets
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
└── rename_app.sh

13 directories, 11 files
```