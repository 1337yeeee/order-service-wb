# Билдер
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /order-service ./cmd/server

# Финальный образ
FROM alpine:3.18
WORKDIR /app
COPY --from=builder /order-service .
EXPOSE 8080
CMD ["./order-service"]
