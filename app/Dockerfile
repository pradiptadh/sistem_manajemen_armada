FROM golang:1.21-alpine

WORKDIR /app

COPY . .

RUN apk add --no-cache git
RUN go mod init pradipta-app && go mod tidy
RUN go build -o main .

CMD ["./main"]