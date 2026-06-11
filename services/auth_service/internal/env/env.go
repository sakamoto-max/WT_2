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

// func Validate() {

// 	_, ok := os.LookupEnv("PG_USER")
// 	if !ok {
// 		log.Fatalf("unable to find env PG_USER")
// 	}
// 	_, ok = os.LookupEnv("PG_PASS")
// 	if !ok {
// 		log.Fatalf("unable to find env PG_PASS")
// 	}
// 	_, ok = os.LookupEnv("PG_HOST")
// 	if !ok {
// 		log.Fatalf("unable to find env PG_HOST")
// 	}
// 	_, ok = os.LookupEnv("PG_ADDRESS")
// 	if !ok {
// 		log.Fatalf("unable to find env PG_ADDRESS")
// 	}
// 	_, ok = os.LookupEnv("PG_DATABASE_NAME")
// 	if !ok {
// 		log.Fatalf("unable to find env PG_DATABASE_NAME")
// 	}
// 	_, ok = os.LookupEnv("PG_SSL_MODE")
// 	if !ok {
// 		log.Fatalf("unable to find env PG_SSL_MODE")
// 	}
// 	_, ok = os.LookupEnv("REDIS_ADDR")
// 	if !ok {
// 		log.Fatalf("unable to find env REDIS_ADDR")
// 	}
// 	_, ok = os.LookupEnv("REDIS_PASS")
// 	if !ok {
// 		log.Fatalf("unable to find env REDIS_PASS")
// 	}
// 	_, ok = os.LookupEnv("REDIS_DB")
// 	if !ok {
// 		log.Fatalf("unable to find env REDIS_DB")
// 	}
// 	_, ok = os.LookupEnv("SECRET_KEY")
// 	if !ok {
// 		log.Fatalf("unable to find env SECRET_KEY")
// 	}
// 	_, ok = os.LookupEnv("GRPC_SERVER_ADDR")
// 	if !ok {
// 		log.Fatalf("unable to find env GRPC_SERVER_ADDR")
// 	}
// 	_, ok = os.LookupEnv("SERVICE_NAME")
// 	if !ok {
// 		log.Fatalf("unable to find env SERVICE_NAME")
// 	}
// }

// # docker run -p 6001:6001 -e POSTGRES_CONN="postgresql://postgres:root@host.docker.internal:5432/auth?sslmode=disable" -e REDIS_ADDR="6379" -e REDIS_DB="0" -e SERVICE_NAME="auth_service" -e REDIS_PASS="" -e SECRET_KEY="asdfghjklazsxdc" -e GRPC_SERVER_ADDR="6001" -it auth_service
