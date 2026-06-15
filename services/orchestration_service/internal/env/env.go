package env

import (
	"log"
	"github.com/joho/godotenv"
)

func Load(filePath string) {
	err := godotenv.Load(filePath)
	if err != nil {
		log.Fatalf("error loading the env file : %v", err)
	}
}
