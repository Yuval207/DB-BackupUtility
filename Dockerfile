# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache make

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN make build

# Runtime stage
FROM alpine:latest

WORKDIR /app

# Install runtime dependencies required by the CLI
RUN apk add --no-cache \
    mysql-client \
    postgresql-client \
    mongodb-tools \
    nodejs \
    npm \
    ca-certificates \
    tzdata

# Copy the binary from the builder stage
COPY --from=builder /app/dbbackup /usr/local/bin/dbbackup

# Set the entrypoint to the CLI
ENTRYPOINT ["dbbackup"]
