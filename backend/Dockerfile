# Multi-stage build for Go backend
FROM golang:1.21-alpine AS builder

# Set working directory
WORKDIR /app

# Install build dependencies
RUN apk --no-cache add git ca-certificates

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy the binary and config files
COPY --from=builder /app/main .
COPY --from=builder /app/.env.yaml.example ./.env.yaml.example

# Expose port
EXPOSE 8080

# Run the binary
CMD ["./main"]