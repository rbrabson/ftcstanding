# Build stage
FROM golang:1.24.0-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git

# Set working directory
WORKDIR /build

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application for Linux AMD64
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -o rank \
    ./cmd/ftc/main.go

# Runtime stage
FROM alpine:latest

# Install ca-certificates for HTTPS connections
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /build/rank .

# Change ownership to non-root user
RUN chown -R 65532:65532 /app

# Switch to non-root user
USER 65532:65532

# Expose any ports if needed (adjust as necessary)
# EXPOSE 8080

# Set environment variables
ENV DATA_SOURCE_NAME=""

# Run the application
CMD ["./rank"]
