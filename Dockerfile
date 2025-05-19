# Start from the official Golang image for building
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install git (required for go mod)
RUN apk add --no-cache git

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the Go app
RUN go build -o dora-handler .

# Start a minimal image for running
FROM alpine:latest
WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/dora-handler .
# Copy config example (optional)
COPY .image-handler.yaml .

# Expose the default port
EXPOSE 8080

# Set environment variable for config file path (optional, can be overridden)
ENV CONFIG=./.image-handler.yaml

# Run the binary
CMD ["./dora-handler", "--config", "/app/.image-handler.yaml"]
