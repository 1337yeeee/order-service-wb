FROM golang:1.24.2-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /order-service ./cmd/server

FROM alpine:3.18
WORKDIR /app
COPY --from=builder /order-service .
EXPOSE 8080
CMD ["./order-service"]
