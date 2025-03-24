# Build stage
FROM golang:1.22-alpine AS builder
WORKDIR /app

# Install dependencies
RUN apk add --no-cache git make

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN go build -o main ./cmd/api

# Run stage
FROM alpine:latest
WORKDIR /app

# Install dependencies
RUN apk add --no-cache ca-certificates tzdata

# Copy binary from build stage
COPY --from=builder /app/main .

# Copy config
COPY --from=builder /app/config/config.yaml ./config/

# Create logs directory
RUN mkdir -p logs

# Expose port
EXPOSE 8080

# Run the application
CMD ["/app/main"]