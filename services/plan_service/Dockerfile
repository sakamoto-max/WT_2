FROM golang:alpine3.21 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main ./cmd/server
RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -o consumer ./cmd/consumer

# ----------------------------------------------------------------------

FROM alpine:3.21

WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/consumer .

# CMD [ "./main" ]