package env

import (
	"log"
	"github.com/joho/godotenv"
)

func Load(fileName string) {
	err := godotenv.Load(fileName)
	if err != nil {
		log.Fatalf("error loading the env file : %v", err)
	}
}


// # docker run -p 6001:6001 -e POSTGRES_CONN="postgresql://postgres:root@host.docker.internal:5432/auth?sslmode=disable" -e REDIS_ADDR="6379" -e REDIS_DB="0" -e SERVICE_NAME="auth_service" -e REDIS_PASS="" -e SECRET_KEY="asdfghjklazsxdc" -e GRPC_SERVER_ADDR="6001" -it auth_service
