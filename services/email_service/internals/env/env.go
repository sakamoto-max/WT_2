package env

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func Load(filePath string) {
	err := godotenv.Load(filePath)
	if err != nil {
		log.Fatalf("failed to load the env file : %v", err)
	}
}

func LookUp() {
	_, ok := os.LookupEnv("MQ_URL")
	if !ok {
		log.Fatalf("unable to find MQ_URL")
	}
}