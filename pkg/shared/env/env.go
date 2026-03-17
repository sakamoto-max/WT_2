package env

import (
	"log"
	"os"
	"github.com/joho/godotenv"
)


func Load() {
	err := godotenv.Load()
	if err != nil{
		log.Fatalf("error loading the env file : %v", err)
	}

	lookup()	
}
func LoadGateway() {
	err := godotenv.Load()
	if err != nil{
		log.Fatalf("error loading the env file : %v", err)
	}

	lookUpForGateway()
}

// -- POSTGRES_CONN
// -- REDIS_ADDR
// -- REDIS_PASS
// -- REDIS_DB
// -- HTTP_SERVER_ADDR
// -- GRPC_SERVER_ADDR

func lookup() {
	_, ok := os.LookupEnv("POSTGRES_CONN")
	if !ok {
		log.Fatalf("POSTGRES_CONN env not found")
	}
	_, ok = os.LookupEnv("REDIS_ADDR")
	if !ok {
		log.Fatalf("REDIS_ADDR env not found")
	}
	_, ok = os.LookupEnv("REDIS_PASS")
	if !ok {
		log.Fatalf("REDIS_PASS env not found")
	}
	_, ok = os.LookupEnv("REDIS_DB")
	if !ok {
		log.Fatalf("REDIS_DB env not found")
	}
	_, ok = os.LookupEnv("HTTP_SERVER_ADDR")
	if !ok {
		log.Fatalf("HTTP_SERVER_ADDR env not found")
	}
	_, ok = os.LookupEnv("GRPC_SERVER_ADDR")
	if !ok {
		log.Fatalf("GRPC_SERVER_ADDR env not found")
	}
	_, ok = os.LookupEnv("SECRET_KEY")
	if !ok {
		log.Fatalf("SECRET_KEY env not found")
	}
}
func lookUpForGateway() {
	_, ok := os.LookupEnv("HTTP_SERVER_ADDR")
	if !ok {
		log.Fatalf("HTTP_SERVER_ADDR env not found")
	}
	_, ok = os.LookupEnv("AUTH_GRPC_CLIENT_ADDR")
	if !ok {
		log.Fatalf("AUTH_GRPC_CLIENT_ADDR env not found")
	}
	_, ok = os.LookupEnv("PLAN_GRPC_CLIENT_ADDR")
	if !ok {
		log.Fatalf("PLAN_GRPC_CLIENT_ADDR env not found")
	}
	_, ok = os.LookupEnv("EXER_GRPC_CLIENT_ADDR")
	if !ok {
		log.Fatalf("EXER_GRPC_CLIENT_ADDR env not found")
	}
	_, ok = os.LookupEnv("TRACKER_GRPC_CLIENT_ADDR")
	if !ok {
		log.Fatalf("TRACKER_GRPC_CLIENT_ADDR env not found")
	}
}
