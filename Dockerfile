FROM golang:1.22.1 AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download 

COPY . .

RUN go build -o main

FROM ubuntu:latest

RUN apt-get update && apt-get install -y git

COPY --from=builder /app/main .

CMD ["./main"]
