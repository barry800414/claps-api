FROM golang:1.19.3-alpine

RUN     mkdir -p /app
WORKDIR /app

COPY go.mod go.mod
COPY go.sum go.sum

RUN go mod download
RUN go build -o app

ENTRYPOINT  ["./app"]