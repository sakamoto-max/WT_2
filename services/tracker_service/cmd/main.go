package main

import (
	"os"
	env "wt/pkg/shared/env"
)

func main() {

	env.Load()

	httpSer := NewhttpServer(os.Getenv("HTTP_SERVER_ADDR"))
	httpSer.Run()
}
