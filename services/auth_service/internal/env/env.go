package env

import (
	"log"
	"os"

	// "github.com/joho/godotenv"
)

// func Load(fileName string) {
// 	err := godotenv.Load(fileName)
// 	if err != nil {
// 		log.Fatalf("error loading the env file : %v", err)
// 	}

// 	lookup()
// }

func Validate() {
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

// # docker run -p 6001:6001 -e POSTGRES_CONN="postgresql://postgres:root@host.docker.internal:5432/auth?sslmode=disable" -e REDIS_ADDR="6379" -e REDIS_DB="0" -e SERVICE_NAME="auth_service" -e REDIS_PASS="" -e SECRET_KEY="asdfghjklazsxdc" -e GRPC_SERVER_ADDR="6001" -it auth_service
