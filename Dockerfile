# Stage 1. Build
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# Static Compilation to evade C library dependencies
RUN CGO_ENABLED=0 GOOS=linux go build -o main cmd/server/main.go

# 2. Final light image
FROM alpine:latest
# Security
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
# Expose the gRPC port
EXPOSE 50051
CMD ["./main"]


