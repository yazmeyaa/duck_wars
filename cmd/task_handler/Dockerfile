FROM golang:1.22.11
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o task_handler ./cmd/task_handler
CMD ["./task_handler"]
