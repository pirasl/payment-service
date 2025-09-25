FROM golang:1.25.1-alpine AS builder

WORKDIR /app

# Copy go mod files first for better layer caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application from cmd/api directory
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o go-app ./cmd/api

FROM alpine:latest

# Install ca-certificates for HTTPS requests (optional but recommended)
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the built binary from the builder stage
COPY --from=builder /app/go-app .

EXPOSE 8080

CMD ["./go-app"]