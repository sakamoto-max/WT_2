package main

import "github.com/go-chi/chi/v5"

func main() {
	// when gets hit with a route -> make grpc call to the particular service
	// should be running in a server
	//

	// run a http server at 5000
	// create clients for all the services
	//

	// r := chi.NewRouter()

	// r.Post("/wt")

}

// routes :
// "/wt/" +

// POST "/user/signup"
// POST "/user/login"
// POST "/user/logout"
// POST "/user/refreshtoken"
// PUT "/user/changepass"
// PUT "/user/changeemail"

// GET "/exercise"
// GET "/exercise/single"
// POST "/exercie"
// DELETE "/exercise"
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
