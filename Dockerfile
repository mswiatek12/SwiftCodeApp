# Use Golang official image to build the app

FROM golang:alpine

# Install required dependencies
RUN apk update && apk add --no-cache git bash build-base

# Set up working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum first to leverage Docker cache
COPY go.mod go.sum ./

# Install Go dependecies
RUN go mod tidy

# Copy the source code into the container
COPY . .

# Build the Go application
RUN go build -o /app/main

# Expose port
EXPOSE 8080

# Set the entrypoint to the compiled Go binary
CMD ["/app/main"]