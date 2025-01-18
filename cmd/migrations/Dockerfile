FROM golang:1.22.11 as builder
WORKDIR /app
RUN go install github.com/pressly/goose/v3/cmd/goose@latest
COPY cmd/migrations/migrations /migrations
WORKDIR /migrations
