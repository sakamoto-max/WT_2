package routes

import (
	"tracker_service/internal/handlers"

	"github.com/go-chi/chi/v5"
	middleware "wt/pkg/middleware"
	chimiddleware "github.com/go-chi/chi/middleware"
)

// think you have a front end

// route to start an empty workout
// route to start a plan workout
// route to end a workout

func Router(h *handlers.Handler) *chi.Mux {
	r := chi.NewRouter()


	r.Use(chimiddleware.Logger)
	r.Use(middleware.JwtMiddleware)
	
	r.Post("/workout/empty", h.StartEmptyWorkout)
	r.Post("/workout", h.StartWorkoutWithPlan)
	r.Post("/workout/end", h.EndWorkout)

	return r
}



// with_out plan :

// {
// 	"plan_id" : 1,
// 	"workout" : [{
// 		"exercise_id" : 23,
// 		"tracker"  : [
// 			{
// 				"set_number" : 1,
// 				"weight" : 20,
// 				"reps" : 10
// 			},
// 			{
// 				"set_number" : 2,
// 				"weight" : 20,
// 				"reps" : 9
// 			}
// 		]
// 	},{
// 		"exercise_id" : 13,
// 		"tracker"  : [
// 			{
// 				"set_number" : 1,
// 				"weight" : 20,
// 				"reps" : 10
// 			},
// 			{
// 				"set_number" : 2,
// 				"weight" : 20,
// 				"reps" : 9
// 			}
// 		]
// 	},{
// 		"exercise_id" : 30,
// 		"tracker"  : [
// 			{
// 				"set_number" : 1,
// 				"weight" : 20,
// 				"reps" : 10
// 			},
// 			{
// 				"set_number" : 2,
// 				"weight" : 20,
// 				"reps" : 9
// 			}
// 		]
// 	},{
// 		"exercise_id" : 15,
// 		"tracker"  : [
// 			{
// 				"set_number" : 1,
// 				"weight" : 20,
// 				"reps" : 10
// 			},
// 			{
// 				"set_number" : 2,
// 				"weight" : 20,
// 				"reps" : 9
// 			}
// 		]
// 	}
// 	]
// }

// tables :
//
