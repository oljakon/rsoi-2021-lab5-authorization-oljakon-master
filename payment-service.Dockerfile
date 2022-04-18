FROM golang:1.17 AS builder

ENV GO111MODULE=on

WORKDIR /app

COPY . .

RUN go mod tidy

RUN go build -o payment ./src/payment-service/cmd/payment-service

ENV PORT=8082
EXPOSE ${PORT}

ENTRYPOINT ["./payment"]