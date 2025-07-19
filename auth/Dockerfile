############################
# 1. Build Stage
############################
FROM golang:1.24-alpine AS builder

# Install necessary tools
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Set Go environment
ENV CGO_ENABLED=0 \
  GOOS=linux \
  GOARCH=amd64

# Copy go.mod and go.sum first (for caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire source code
COPY . .

# Build the Go binary (static)
RUN go build -o server ./cmd/server

############################
# 2. Runtime Stage (scratch)
############################
FROM scratch AS prod

# Copy built binary only
COPY --from=builder /app/server /app/server

# Create non-root user
USER 1001

# Set environment variables
ENV ENV=prod
ENV HTTP_PORT=8080
ENV GRPC_PORT=50051

# Expose ports
EXPOSE 8080
EXPOSE 50051

ENTRYPOINT ["/app/server"]
