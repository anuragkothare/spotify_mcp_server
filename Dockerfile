# Build stage
FROM golang:1.23-alpine AS builder

# Set build arguments
ARG TARGETPLATFORM
ARG BUILDPLATFORM
ARG TARGETOS
ARG TARGETARCH

# Add metadata
LABEL maintainer="anuragkothare7x@gmail.com"
LABEL description="Spotify MCP Server"
LABEL version="1.0.0"

# Set working directory
WORKDIR /build

# Install build dependencies
RUN apk add --no-cache \
    git \
    ca-certificates \
    tzdata \
    upx

# Copy dependency files first (better caching)
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download && go mod verify

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 \
    GOOS=${TARGETOS:-linux} \
    GOARCH=${TARGETARCH:-amd64} \
    go build \
    -ldflags="-w -s -X main.version=1.0.0 -X main.buildTime=$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
    -a \
    -installsuffix cgo \
    -o spotify-mcp-server \
    ./cmd/server

# Compress binary (optional, reduces size by ~30%)
RUN upx --best --lzma spotify-mcp-server

# =================
# PRODUCTION STAGE
# =================
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add \
    ca-certificates \
    tzdata \
    curl \
    && rm -rf /var/cache/apk/*

# Create non-root user
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Set working directory
WORKDIR /app

# Copy binary and configs from builder
COPY --from=builder /build/spotify-mcp-server .
COPY --from=builder /build/configs ./configs

# Create necessary directories
RUN mkdir -p /app/logs /app/tmp

# Set proper ownership
RUN chown -R appuser:appgroup /app

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 8080

# Add health check
HEALTHCHECK --interval=30s \
    --timeout=10s \
    --start-period=5s \
    --retries=3 \
    CMD curl -f http://localhost:8080/health || exit 1

# Set entrypoint
ENTRYPOINT ["./spotify-mcp-server"]