FROM golang:1.26.2-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build \
    -trimpath \
    -ldflags="-s -w" \
    -o /app/bin/seller-service \
    ./cmd/seller-service


FROM alpine:latest

WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY --from=builder /app/bin/seller-service /app/seller-service
COPY config /app/config

ENV CONFIG_PATH=/app/config/local.yaml

EXPOSE 44045

ENTRYPOINT ["/app/seller-service"]