FROM golang:1.24.2 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o user_service ./cmd/user_service/main.go


FROM alpine:latest

COPY --from=builder /app/user_service /usr/local/bin/user_service
COPY config.yaml /app/config.yaml

COPY scripts/wait-db-connect.sh /app/wait-db-connect.sh
RUN chmod +x /app/wait-db-connect.sh


EXPOSE 50051

ENTRYPOINT ["/app/wait-db-connect.sh", "db", "5432"]
CMD [ "/usr/local/bin/user_service" ]
