package env

import (
	"log"
	"os"
)


func LookUp() {
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
}
