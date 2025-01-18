FROM golang:1.22.11
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o telegram_bot ./cmd/telegram_bot
CMD ["./telegram_bot"]
