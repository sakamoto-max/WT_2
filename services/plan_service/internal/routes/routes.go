package routes

import (
	handler "plan_service/internal/handlers"
	"plan_service/internal/middleware"

	chimiddleware "github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

func Router(h *handler.Handler) *chi.Mux {
	r := chi.NewRouter()

	r.Use(chimiddleware.Logger)
	r.Get("/health", h.CheckHealth)

	// create a plan
	r.With(middleware.JwtMiddleware).Post("/plan/create", h.CreatePlan) // done
	// update a plan
	r.Put("/plan/{planName}", h.UpdateThePlan)
	// delete a plan
	r.Delete("/plan/{planName}", h.DeleteAPlan)
	// get all plans
	r.With(middleware.JwtMiddleware).Get("/plan", h.GetAllPlans)
	// get one plan
	r.With(middleware.JwtMiddleware).Get("/plan/{planName}", h.GetPLanByName) // done

	return r
}
