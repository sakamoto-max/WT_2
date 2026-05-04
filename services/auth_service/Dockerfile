FROM golang:1.26-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main ./cmd/server/

# RUN go build -o main ./cmd/server/
# -------------------------------------------------------

FROM alpine:3.21

WORKDIR /app

COPY --from=builder /app/main .

EXPOSE 6001

CMD ["./main"]

# docker run -p 6001:6001 -e POSTGRES_CONN="postgresql://docker_postgres:root@localhost:5432/auth" -e REDIS_DB="" -e REDIS_ADDR="" -e REDIS_DB="" -e SERVICE_NAME="" -e