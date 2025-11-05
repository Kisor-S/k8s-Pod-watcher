# -------Build Stage-------
FROM golang:1.21-alpine AS builder

# Install git (required for go modules) and ca-certificates
RUN apk update && apk add --no-cache git ca-certificates

# Set the working directory inside the container
WORKDIR /app

# Enable Go modules and set up environment
ENV CGO_ENABLED=0 GOFLAGS=-mod=mod

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy rest of the application source code
COPY . .

# Build static binary
RUN go build -o /out/pod-watcher ./main.go

# -------Final Stage-------
FROM alpine:latest
RUN apk --no-cache add ca-certificates

# Create a non-root user (optional but recommended for security)
RUN addgroup -S app && adduser -S app -G app
USER app

COPY --from=builder /out/pod-watcher /usr/local/bin/pod-watcher

# Default args: use kubeconfig or in-cluster config if not provided
ENTRYPOINT ["/usr/local/bin/pod-watcher"]


