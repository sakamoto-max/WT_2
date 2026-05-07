package env

import (
	"log"
	"os"

	// "github.com/joho/godotenv"
	"github.com/joho/godotenv"
)

func Load(fileName string) {
	err := godotenv.Load(fileName)
	if err != nil {
		log.Fatalf("error loading the env file : %v", err)
	}

	LookUp()
}

func LookUp() {
	_, ok := os.LookupEnv("AUTH_POSTGRES_CONN")
	if !ok {
		log.Fatalf("unable to find env AUTH_POSTGRES_CONN")
	}
	_, ok = os.LookupEnv("REDIS_ADDR")
	if !ok {
		log.Fatalf("unable to find env REDIS_ADDR")
	}
	os.LookupEnv("REDIS_PASS")
	if !ok {
		log.Fatalf("unable to find env REDIS_PASS")
	}
	os.LookupEnv("REDIS_DB")
	if !ok {
		log.Fatalf("unable to find env REDIS_DB")
	}
	os.LookupEnv("TRACKER_POSTGRES_CONN")
	if !ok {
		log.Fatalf("unable to find env SECRET_KEY")
	}
	os.LookupEnv("MQ_URL")
	if !ok {
		log.Fatalf("unable to find env MQ_URL")
	}
}
