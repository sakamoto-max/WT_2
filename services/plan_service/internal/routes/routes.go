package routes

// import (
// 	handler "plan_service/internal/handlers"
// 	// "plan_service/internal/middleware"
// 	middleware "wt/pkg/middleware"

// 	chimiddleware "github.com/go-chi/chi/middleware"
// 	"github.com/go-chi/chi/v5"
// )

// func Router(h *handler.Handler) *chi.Mux {
// 	r := chi.NewRouter()

// 	r.Use(chimiddleware.Logger)
// 	r.Get("/health", h.CheckHealth)

// 	// create a plan
// 	r.With(middleware.JwtMiddleware).Post("/plan/create", h.CreatePlan) // done
// 	// add exercises to an existing plan
// 	r.With(middleware.JwtMiddleware).Put("/plan/exercises", h.AddExercisesToPlan)
// 	// delete exercises from an existing plan
// 	r.With(middleware.JwtMiddleware).Delete("/plan/exercises", h.DeleteExerciseFromPlan)
// 	// delete a plan
// 	r.With(middleware.JwtMiddleware).Delete("/plan", h.DeletePlan)
// 	// get all plans
// 	r.With(middleware.JwtMiddleware).Get("/plan", h.GetAllPlans)
// 	// get one plan
// 	r.With(middleware.JwtMiddleware).Get("/plan/oneplan", h.GetPLanByName) // done

// 	return r
// }
