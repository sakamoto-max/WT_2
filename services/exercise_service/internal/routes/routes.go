package routes

import (
	"exercise_service/internal/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func Routes(h *handlers.Handler) *chi.Mux {
	// create exercise
	// get all exercises
	// update exercises
	// delete exercise

	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Get("/exercise", h.GetAllExercises)
	r.Get("/exercise/{exerciseName}", h.GetExerciseByName)
	r.Post("/exercise", h.CreateExercise)
	r.Delete("/exercise/{exerciseName}", h.DeleteExecise)

	return r
}
