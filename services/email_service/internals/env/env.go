package env

import (
	"log"
	"os"
)

func LookUp() {
	_, ok := os.LookupEnv("MQ_URL")
	if !ok {
		log.Fatalf("unable to find MQ_URL")
	}
}