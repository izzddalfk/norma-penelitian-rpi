FROM --platform=$BUILDPLATFORM golang:1.20-alpine3.18 AS builder

WORKDIR /norma/penelitian-rpi

COPY go.mod .
COPY go.sum .
COPY . .

ARG TARGETARCH

RUN go mod download
RUN GOOS=linux GOARCH=$TARGETARCH go build -a -o app ./cmd/umkm/main.go

FROM alpine:3.18

WORKDIR /norma/penelitian-rpi

COPY --from=builder /norma/penelitian-rpi/app .
ENTRYPOINT [ "/norma/penelitian-rpi/app" ]