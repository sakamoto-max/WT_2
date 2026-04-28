package main

import (
	"os"

	"github.com/sakamoto-max/wt_2/api_gateway/env"
)

func main() {
	env.Load("../.env")
	
	
	httpSer := NewHttpServer(os.Getenv("HTTP_SERVER_ADDR"))
	httpSer.Run()
}

// r.Post("/wt/user/login", handler.)

// routes :
// "/wt/" +
// PUT "/user/changepass"
// PUT "/user/changeemail"

// PUT "/exercise/{exercisename}"

// GET "/plan/health"
// POST "/plan/create"
// PUT "/plan/exercises"
// DELETE "/plan/exercises"
// GET "/plan"
// GET "/plan/oneplan"

// POST "/workout/empty"
// POST "/workout"
// POST "/workout/end"

// analytics

// "/analytics/bodypart/{bodypartname}"
// "/analytics/plan/{planname}"
// "/analytics/exercise/{exercisename}"
// "/anlaytics"
