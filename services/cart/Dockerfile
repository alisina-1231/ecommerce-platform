FROM golang:1.24.4 AS builder

WORKDIR /app

COPY go.mod ./

COPY cmd ./cmd

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o cart-service ./cmd/server

FROM alpine:3.20

RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

COPY --from=builder /app/cart-service .

RUN chown appuser:appgroup /app/cart-service

USER appuser

EXPOSE 8080

ENV PORT=8080

CMD ["./cart-service"]