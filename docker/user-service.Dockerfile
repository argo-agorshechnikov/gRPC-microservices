FROM golang:1.24.2 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o user_service ./cmd/user_service/main.go


FROM alpine:latest

COPY --from=builder /app/user_service /usr/local/bin/user_service
COPY config.yaml /app/config.yaml


EXPOSE 50051

ENTRYPOINT ["/usr/local/bin/user_service"]
