# Build Stage
FROM golang:1.24.1 AS builder

# Set working directory
WORKDIR /app

# Set Go cache and mod cache to writable directories
ENV GOCACHE=/app/.cache/go-build
ENV GOMODCACHE=/app/.cache/go-mod

# Copy go.mod and go.sum first to download dependencies
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the binary
RUN go build -o /bin/swiftcode_app ./main

# Run Stage (smaller final image)
FROM ubuntu:latest

# Install CA certificates to enable TLS verification
RUN apt-get update && apt-get install -y ca-certificates wget curl postgresql-client && rm -rf /var/lib/apt/lists/*

# Install Go in the final image (to run tests, if needed)
RUN wget https://go.dev/dl/go1.24.1.linux-amd64.tar.gz && \
    tar -C /usr/local -xvzf go1.24.1.linux-amd64.tar.gz && \
    rm go1.24.1.linux-amd64.tar.gz

# Set Go binary path
ENV PATH="/usr/local/go/bin:${PATH}"

# Copy configuration files and the built binary from the builder stage
COPY --from=builder /bin/swiftcode_app /app/swiftcode_app
COPY config/.env /app/config/.env
COPY credentials.json /app/credentials.json

# Copy the entire application code (handlers, models, services, etc.)
COPY handlers/ /app/handlers
COPY main/ /app/main
COPY repository/ /app/repository
COPY models/ /app/models
COPY service/ /app/service

# Copy tests directory (if you need it to run tests in container)
COPY tests/ /app/tests/

# Expose the application port
EXPOSE 8080

# Command to run the application
CMD ["/app/swiftcode_app"]
