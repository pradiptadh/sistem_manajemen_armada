FROM golang:1.23-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN  go mod tidy && GOOS=linux GOARCH=amd64 go build -v -o main . && chmod +x main

CMD ["./main"]