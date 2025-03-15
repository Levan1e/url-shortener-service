ARG GO_VERSION=1.23
ARG ALPINE_VERSION=3.21

FROM golang:${GO_VERSION}-alpine${ALPINE_VERSION} AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o url-shortener-service cmd/shortener/main.go

ARG ALPINE_VERSION=3.21
FROM alpine:${ALPINE_VERSION} AS deploy

WORKDIR /app

COPY --from=builder /app/url-shortener-service .
COPY internal/config/config.yaml /app/internal/config/config.yaml

ENV STORAGE=memory

CMD ["sh", "-c", "./url-shortener-service --storage=$STORAGE"]

