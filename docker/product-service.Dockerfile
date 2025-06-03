FROM golang:1.24.2 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY ./cmd/product_service ./cmd/product_service
COPY ./internal ./internal
COPY ./api ./api
COPY ./pkg ./pkg

RUN CGO_ENABLED=0 GOOS=linux go build -o product_service ./cmd/product_service/main.go

FROM alpine:latest

COPY --from=builder /app/product_service /usr/local/bin/product_service

EXPOSE 50052

ENTRYPOINT [ "/usr/local/bin/product_service" ]