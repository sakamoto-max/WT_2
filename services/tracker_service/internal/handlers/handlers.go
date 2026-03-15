package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
	"tracker_service/internal/models"
	"tracker_service/internal/services"
	"tracker_service/internal/user"
	"tracker_service/internal/utils"
	token "wt/pkg/shared"
)

type Handler struct {
	service *services.Service
}

func NewHandler(s *services.Service) *Handler {
	return &Handler{service: s}
}

// req ;
// empty

// resp :
// {
// 	"message" : 'empty workout has started'
// }

func (h *Handler) StartEmptyWorkout(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*3)
	defer cancel()

	t := token.JwtToken{}

	claims, ok := t.GetClaimsFromContext(ctx)
	if !ok {

	}

	err := h.service.StartEmptyWorkoutSer(ctx, claims.UserId)
	if err != nil {
		utils.BadReqErrorWriter(w, err.Error())
		return
	}

	resp := models.GeneralResp{}
	resp.Message = "an empty workout has started"

	utils.CreatedResponseWriter(w, resp)
}

// req :
// {
// 	"plan_name" : "plan"
// }

// resp :
// {
// 	"message" : 'plan workout has started'
// }

func (h *Handler) StartWorkoutWithPlan(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*3)
	defer cancel()

	var planName user.PlanName

	json.NewDecoder(r.Body).Decode(&planName)

	validationErr, errOccured := planName.Validate()
	if errOccured {
		utils.ValidationErrWriter(w, validationErr)
		return
	}

	t := token.JwtToken{}
	claims, ok := t.GetClaimsFromContext(ctx)
	if !ok {

	}
	resp, err := h.service.StartWorkoutWithPlanSer(ctx, claims.UserId, planName.PlanName)
	if err != nil {
		utils.BadReqErrorWriter(w, err.Error())
		return
	}
	utils.OkResponseWriter(w, resp)

}

// req :
// {
// 	"workout" : [{
// 		"exercise_id" : 23,
// 		"tracker"  : [
// 			{
// 				"weight" : 20,
// 				"reps" : 10
// 			},
// 			{
// 				"weight" : 20,
// 				"reps" : 9
// 			}
// 		]
// 	},{
// 		"exercise_id" : 13,
// 		"tracker"  : [
// 			{
// 				"weight" : 20,
// 				"reps" : 10
// 			},
// 			{
// 				"weight" : 20,
// 				"reps" : 9
// 			}
// 		]
// 	},{
// 		"exercise_id" : 30,
// 		"tracker"  : [
// 			{
// 				"weight" : 20,
// 				"reps" : 10
// 			},
// 			{
// 				"weight" : 20,
// 				"reps" : 9
// 			}
// 		]
// 	},{
// 		"exercise_id" : 15,
// 		"tracker"  : [
// 			{
// 				"weight" : 20,
// 				"reps" : 10
// 			},
// 			{
// 				"weight" : 20,
// 				"reps" : 9
// 			}
// 		]
// 	}
// 	]
// }

// resp :
// {
// 	"message" : "hurray, workout has ended successfully."
// }

func (h *Handler) EndWorkout(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*3)
	defer cancel()
	var req user.Tracker

	json.NewDecoder(r.Body).Decode(&req)

	validationErr, errOccured := req.Validate()
	if errOccured {
		utils.ValidationErrWriter(w, validationErr)
		return
	}

	t := token.JwtToken{}
	claims, ok := t.GetClaimsFromContext(ctx)
	if !ok {

	}

	err := h.service.EndWorkoutSer(ctx, claims.UserId, &req)
	if err != nil{
		utils.BadReqErrorWriter(w, err.Error())
		return
	}

	resp := map[string]string{
		"message" : "workout ended successfully",
	}

	utils.OkResponseWriter(w, resp)
}
