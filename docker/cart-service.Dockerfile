FROM golang:1.24.2 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o cart_service ./cmd/cart_service/main.go

FROM alpine:latest

COPY --from=builder /app/cart_service /usr/local/bin/cart_service
COPY config.yaml /app/config.yaml

COPY scripts/wait-db-connect.sh /app/wait-db-connect.sh
RUN chmod +x /app/wait-db-connect.sh

EXPOSE 50053


ENTRYPOINT ["/app/wait-db-connect.sh", "db", "5432"]
CMD [ "/usr/local/bin/cart_service" ]

