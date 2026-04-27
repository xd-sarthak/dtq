# Build stage
FROM golang:1.24-bookworm AS builder

WORKDIR /app

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Build the binary — CGO_ENABLED=1 to link glibc DNS resolver
COPY . .
RUN go build -o /app/server ./cmd/server

# Runtime stage
FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY --from=builder /app/server .

EXPOSE 8080

ENV GODEBUG=netdns=cgo
CMD ["./server"]
