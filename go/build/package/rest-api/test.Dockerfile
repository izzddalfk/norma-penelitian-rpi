FROM golang:1.20-alpine3.18

WORKDIR /norma/penelitian-rpi/go

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY ./internal .