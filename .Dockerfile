FROM golang:1.24-alpine
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o main ./cmd/server
CMD ["./main"]
