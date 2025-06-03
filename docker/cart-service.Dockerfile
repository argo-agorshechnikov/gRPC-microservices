FROM golang:1.24.2 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY ./cmd/cart_service ./cmd/cart_service
COPY ./internal ./internal
COPY ./api ./api
COPY ./pkg ./pkg

RUN CGO_ENABLED=0 GOOS=linux go build -o cart_service ./cmd/cart_service/main.go

FROM alpine:latest

COPY --from=builder /app/cart_service /usr/local/bin/cart_service

EXPOSE 50053

ENTRYPOINT [ "/usr/local/bin/cart_service" ]

