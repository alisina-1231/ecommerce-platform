FROM golang:1.26.3 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY cmd ./cmd
COPY internal ./internal

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o checkout-service ./cmd/server

FROM alpine:3.20

RUN addgroup -S appgroup && \
    adduser -S appuser -G appgroup

WORKDIR /app

COPY --from=builder /app/checkout-service .

RUN chown appuser:appgroup /app/checkout-service

USER appuser

EXPOSE 8080

ENV PORT=8080

CMD ["./checkout-service"]