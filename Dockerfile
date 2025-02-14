# Build stage
FROM golang:alpine AS builder
WORKDIR /app
RUN apk add --no-cache git
COPY . .
RUN go mod download
RUN go build -o crypto_bot ./cmd/crypto_bot

# Runtime stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/crypto_bot .
CMD ["./crypto_bot"]
