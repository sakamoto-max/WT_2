FROM go:1.26.2-alpine3.23 AS builder

RUN mkdir app

WORKDIR /app

COPY go.mod go.sum /app/
RUN go mod download

RUN go build -o main

