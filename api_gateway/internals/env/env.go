package env

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	// "github.com/joho/godotenv"
)

func Load(fileName string) {
	err := godotenv.Load(fileName)
	if err != nil {
		log.Fatalf("error loading the env file : %v", err)
	}
}

func LookUp() {
	_, ok := os.LookupEnv("HTTP_SERVER_ADDR")
	if !ok {
		log.Fatalf("unable to find env HTTP_SERVER_ADDR")
	}
	_, ok = os.LookupEnv("AUTH_GRPC_CLIENT_ADDR")
	if !ok {
		log.Fatalf("unable to find env AUTH_GRPC_CLIENT_ADDR")
	}
	_, ok = os.LookupEnv("PLAN_GRPC_CLIENT_ADDR")
	if !ok {
		log.Fatalf("unable to find env PLAN_GRPC_CLIENT_ADDR")
	}
	_, ok = os.LookupEnv("EXER_GRPC_CLIENT_ADDR")
	if !ok {
		log.Fatalf("unable to find env EXER_GRPC_CLIENT_ADDR")
	}
	_, ok = os.LookupEnv("SECRET_KEY")
	if !ok {
		log.Fatalf("unable to find env SECRET_KEY")
	}
	_, ok = os.LookupEnv("TRACKER_GRPC_CLIENT_ADDR")
	if !ok {
		log.Fatalf("unable to find env TRACKER_GRPC_CLIENT_ADDR")
	}
	_, ok = os.LookupEnv("SERVICE_NAME")
	if !ok {
		log.Fatalf("unable to find env SERVICE_NAME")
	}
}
