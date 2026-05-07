package env

import (
	"log"
	"os"
)


func Lookup() {
	_, ok := os.LookupEnv("POSTGRES_CONN")
	if !ok {
		log.Fatalf("unable to find env POSTGRES_CONN")
	}
	_, ok = os.LookupEnv("REDIS_ADDR")
	if !ok {
		log.Fatalf("unable to find env REDIS_ADDR")
	}
	_, ok = os.LookupEnv("REDIS_PASS")
	if !ok {
		log.Fatalf("unable to find env REDIS_PASS")
	}
	_, ok = os.LookupEnv("REDIS_DB")
	if !ok {
		log.Fatalf("unable to find env REDIS_DB")
	}
	_, ok = os.LookupEnv("SECRET_KEY")
	if !ok {
		log.Fatalf("unable to find env SECRET_KEY")
	}
	_, ok = os.LookupEnv("GRPC_SERVER_ADDR")
	if !ok {
		log.Fatalf("unable to find env GRPC_SERVER_ADDR")
	}
	_, ok = os.LookupEnv("SERVICE_NAME")
	if !ok {
		log.Fatalf("unable to find env SERVICE_NAME")
	}
	_, ok = os.LookupEnv("EXERCISE_GRPC_SERVER_ADDR")
	if !ok {
		log.Fatalf("unable to find env EXERCISE_GRPC_SERVER_ADDR")
	}
	_, ok = os.LookupEnv("MQ_URL")
	if !ok {
		log.Fatalf("unable to find env MQ_URL")
	}
}

//  docker run -p 5000:5000 -e POSTGRES_CONN="postgresql://postgres:root@host.docker.internal:5432/plan?sslmode=disable" -e REDIS_ADDR="host.docker.internal:6379" -e REDIS_PASS="" -e REDIS_DB="1" -e SECRET_KEY="asdfghjklazsxdc" -e GRPC_SERVER_ADDR=":6002" -e SERVICE_NAME="plan_service" -e EXERCISE_GRPC_SERVER_ADDR=":6003" -e MQ_URL="amqp://guest:guest@host.docker.internal:5672/" -it plan_service

// docker run -p 6001:6001 -e POSTGRES_CONN="postgresql://postgres:root@host.docker.internal:5432/auth?sslmode=disable" -e REDIS_ADDR="host.docker.internal:6379" -e REDIS_DB="0" -e SERVICE_NAME="auth_service" -e REDIS_PASS="" -e SECRET_KEY="asdfghjklazsxdc" -e GRPC_SERVER_ADDR="6001" -it auth_service