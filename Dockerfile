FROM golang:1.20.1-alpine3.17 AS builder

WORKDIR /app

COPY . /app

RUN go build -o main .

# Set the command to run the binary
CMD ["/app/main"]
