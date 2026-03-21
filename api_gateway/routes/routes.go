package routes

import (
	"api_gateway/handlers"
	// "net/http"
	"wt/pkg/middleware"

	chimiddleware "github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

func NewRouter(h *handlers.Handler) *chi.Mux {

	r := chi.NewRouter()

	r.Use(chimiddleware.Logger)

	r.Get("/wt/health", h.GetHealth)

	r.Post("/wt/user/signup", h.SignUp)
	r.Post("/wt/user/login", h.Login)
	r.With(middleware.JwtMiddleware).Post("/wt/user/logout", h.Logout)
	r.Post("/wt/user/refresh", h.GetNewAccessToken)
	r.With(middleware.JwtMiddleware).Patch("/wt/user/password", h.ChangePassWord)
	r.With(middleware.JwtMiddleware).Patch("/wt/user/email", h.ChangeEmail)

	r.Get("/wt/exercises", h.GetAllExercises)
	r.Get("/wt/exercises/single", h.GetExerciseByName)
	r.With(middleware.JwtMiddleware).Post("/wt/exercise", h.CreateExercise)
	r.With(middleware.JwtMiddleware).Delete("/wt/exercise", h.DeleteExecise)

	// r.Get("/wt/plan/health", h.CheckHealthPlan)
	r.With(middleware.JwtMiddleware).Post("/wt/plan/create", h.CreatePlan)
	r.With(middleware.JwtMiddleware).Patch("/wt/plan/exercises", h.AddExercisesToPlan)
	r.With(middleware.JwtMiddleware).Delete("/wt/plan/exercises", h.DeleteExerciseFromPlan)
	r.With(middleware.JwtMiddleware).Get("/wt/plan/health", h.CheckHealthPlan)
	r.With(middleware.JwtMiddleware).Get("/wt/plan", h.GetAllPlans)
	r.With(middleware.JwtMiddleware).Get("/wt/plan/oneplan", h.GetPLanByName)
	r.With(middleware.JwtMiddleware).Delete("/wt/plan", h.DeletePlan)

	r.With(middleware.JwtMiddleware).Post("/wt/workout/empty", h.StartEmptyWorkout)
	r.With(middleware.JwtMiddleware).Post("/wt/workout", h.StartWorkoutWithPlan)
	r.With(middleware.JwtMiddleware).Post("/wt/workout/end", h.EndWorkout)

	// http.ListenAndServe(":5000", r)
	return r

}

// func NewRouter() {
// 	r := chi.NewRouter()

// 	r.Post("/wt/user/signup")
// 	r.Post("/wt/user/login")
// 	r.Post("/wt/user/logout")
// 	r.Post("/wt/user/refreshtoken")
// 	r.Patch("/wt/user/changepass")
// 	r.Patch("/wt/user/changeemail")

// 	r.Get("/wt/exercise")
// 	r.

// }

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
