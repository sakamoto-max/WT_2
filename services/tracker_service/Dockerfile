FROM golang:alpine3.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLE=0 GOARCH=amd64 GOOS=linux go build -o main ./cmd/server

# ------------------------------------------------------------------

FROM alpine:3.21

WORKDIR /app

COPY --from=builder /app/main .

EXPOSE 6003

CMD [ "./main" ]


