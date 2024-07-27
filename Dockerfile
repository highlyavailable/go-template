FROM golang:1.20

WORKDIR /goapp

# Download Go modules
COPY goapp/go.mod goapp/go.sum ./
RUN go mod download

# Copy source code
COPY goapp .

# Copy .env to /goapp/.env in the container
COPY .env /goapp/.env

# Build the application and produce a binary in the build folder
RUN go build -o build/goapp ./cmd/goapp

# Expose port if your application listens on a port
EXPOSE 8080

# Make the binary executable
RUN chmod +x ./build/goapp

# Run the application
CMD ["/goapp/build/goapp"]