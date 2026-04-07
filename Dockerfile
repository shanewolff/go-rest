# Build stage
FROM golang:1.26-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/api/main.go

# --- Debug Build ---
FROM builder AS debug-builder
# Install Delve
RUN go install github.com/go-delve/delve/cmd/dlv@latest
# Build with optimizations disabled
RUN CGO_ENABLED=0 GOOS=linux go build -gcflags="all=-N -l" -o main-debug ./cmd/api/main.go

# --- Final Stage ---
FROM alpine:latest AS final

RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
EXPOSE 8080
CMD ["./main"]

# --- Debug Final Stage ---
FROM alpine:latest AS debug
RUN apk --no-cache add ca-certificates libc6-compat
WORKDIR /root/
COPY --from=debug-builder /app/main-debug ./main
COPY --from=debug-builder /go/bin/dlv /usr/local/bin/dlv
EXPOSE 8080 40000
# Run the application via Delve
CMD ["dlv", "--listen=:40000", "--headless=true", "--api-version=2", "--accept-multiclient", "exec", "./main"]
