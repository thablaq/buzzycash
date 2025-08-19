# Stage 1: Build
FROM golang:1.21-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

# Copy go.mod and go.sum first to leverage Docker cache
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the Go binary (main.go is under ./cmd)
RUN go build -o server ./cmd

# Stage 2: Run
FROM alpine:latest
RUN apk add --no-cache ca-certificates

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/server .

EXPOSE 5005

# Run the binary
CMD ["./server"]
