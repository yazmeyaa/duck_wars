FROM golang:1.22.11
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o /godocker ./cmd/task_handler/task_handler.go
EXPOSE 8080
CMD [ "/godocker" ]
