FROM golang:1.17 AS builder

ENV GO111MODULE=on

WORKDIR /app

COPY . .

RUN go mod tidy

RUN go build -o car ./src/car-service/cmd/car-service

ENV PORT=8081
EXPOSE ${PORT}

ENTRYPOINT ["./car"]