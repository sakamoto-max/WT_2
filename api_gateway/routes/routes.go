package routes

import (
	"github.com/sakamoto-max/wt_2/api_gateway/handlers"
	"github.com/sakamoto-max/wt_2-pkg/middleware"
	"github.com/go-chi/chi/v5"
)

func NewRouter(h *handlers.Handler) *chi.Mux {

	r := chi.NewRouter()
	r.Use(middleware.ReqIdGenerator)
	r.Use(middleware.Logger)
	r.Get("/wt/health", h.GetHealth)

	r.Post("/wt/user/signup", h.SignUp)
	r.Post("/wt/user/login", h.Login)
	r.With(middleware.JwtMiddleware).Post("/wt/user/logout", h.Logout)
	r.Post("/wt/user/refresh", h.GetNewAccessToken)
	r.With(middleware.JwtMiddleware).Patch("/wt/user/password", h.ChangePassWord)
	r.With(middleware.JwtMiddleware).Patch("/wt/user/email", h.ChangeEmail)

	r.With(middleware.JwtMiddleware).Get("/wt/exercises", h.GetAllExercises)
	r.With(middleware.JwtMiddleware).Get("/wt/exercises/{exerciseName}", h.GetExerciseByName)
	r.With(middleware.JwtMiddleware).Post("/wt/exercises", h.CreateExercise)
	r.With(middleware.JwtMiddleware).Delete("/wt/exercises", h.DeleteExecise)

	r.With(middleware.JwtMiddleware).Post("/wt/plan/create", h.CreatePlan)
	r.With(middleware.JwtMiddleware).Patch("/wt/plan/exercises", h.AddExercisesToPlan)
	r.With(middleware.JwtMiddleware).Delete("/wt/plan/exercises", h.DeleteExerciseFromPlan)
	r.With(middleware.JwtMiddleware).Get("/wt/plan", h.GetAllPlans)
	r.With(middleware.JwtMiddleware).Get("/wt/plan/oneplan", h.GetPLanByName)
	r.With(middleware.JwtMiddleware).Delete("/wt/plan", h.DeletePlan)

	r.With(middleware.JwtMiddleware).Post("/wt/workout/empty", h.StartEmptyWorkout)
	r.With(middleware.JwtMiddleware).Post("/wt/workout", h.StartWorkoutWithPlan)
	r.With(middleware.JwtMiddleware).Post("/wt/workout/end", h.EndWorkout)
	r.With(middleware.JwtMiddleware).Post("/wt/workout/cancel", h.CancelWorkout)

	return r
}

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
