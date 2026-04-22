package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	exerpb "workout-tracker/proto/shared/exercise"
	"wt/pkg/middleware"
	myerrors "wt/pkg/my_errors"
	"wt/pkg/user"
	"wt/pkg/utils"

	"go.uber.org/zap"
)

func (h *Handler) CreateExercise(w http.ResponseWriter, r *http.Request) {

	logger, err := middleware.GetLogger(r.Context())
	if err != nil {
		utils.InternalServerErr(w, err)
	}
	reqId, err := middleware.GetReqId(r.Context())
	if err != nil {
		utils.InternalServerErr(w, err)
	}
	claims, err := middleware.GetClaims(r.Context())
	if err != nil {
		utils.InternalServerErr(w, err)
	}

	logger.Log.Infow("CREATE_EXERCISE_CALLED", zap.String("REQ_ID", reqId))

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*3)
	defer cancel()

	var userInput user.Exercise

	json.NewDecoder(r.Body).Decode(&userInput)

	validationErr, errOccured := userInput.Validate()
	if errOccured {
		utils.ValidationErrWriter(w, *validationErr)
		return
	}

	in := exerpb.CreateExerciseReq{
		ExerciseName: userInput.Name,
		BodyPart:     userInput.BodyPart,
		Equipment:    userInput.Equipment,
		UserId:       claims.UserId,
	}

	res, err := h.exerClient.CreateExercise(ctx, &in)
	if err != nil {
		myerrors.ErrMatcher(w, err)
		return
	}

	resp := user.CreateExerResp{
		Messsage: fmt.Sprintf("exercise %v has been successfully created ", in.ExerciseName),
		Exercise: user.Exercise{
			Id:        res.Id,
			Name:      in.ExerciseName,
			BodyPart:  in.BodyPart,
			Equipment: in.Equipment,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	utils.CreatedWriter(w, resp)
	logger.Log.Infow("EXERCISE_CREATION_SUCCESSFULL", zap.String("REQ_ID", reqId))
}

func (h *Handler) GetExerciseByName(w http.ResponseWriter, r *http.Request) {

	claims, err := middleware.GetClaims(r.Context())
	if err != nil {
		utils.InternalServerErr(w, err)
	}
	logger, err := middleware.GetLogger(r.Context())
	if err != nil {
		utils.InternalServerErr(w, err)
	}
	reqId, err := middleware.GetReqId(r.Context())
	if err != nil {
		utils.InternalServerErr(w, err)
	}

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*3)
	defer cancel()

	logger.Log.Infow("GET_EXERCISE_CALLED", zap.String("REQ_ID", reqId))

	Name := r.PathValue("exerciseName")

	fmt.Println("making the in")

	in := exerpb.SendExerciseName{
		ExerciseName: Name,
		UserId:       claims.UserId,
	}

	fmt.Println("sending the in")
	res, err := h.exerClient.GetOneExercise(ctx, &in)
	if err != nil {
		myerrors.ErrMatcher(w, err)
		return
	}

	resp := user.Exercise{
		Id:        res.Id,
		Name:      res.Name,
		BodyPart:  res.BodyPart,
		Equipment: res.Equipment,
		CreatedAt: res.CreatedAt.AsTime(),
		UpdatedAt: res.UpdatedAt.AsTime(),
	}

	utils.OkRespWriter(w, resp)
	logger.Log.Infow("GET_EXERCISE_SUCCESSFULL", zap.String("REQ_ID", reqId))
}

func (h *Handler) GetAllExercises(w http.ResponseWriter, r *http.Request) {
	claims, err := middleware.GetClaims(r.Context())
	if err != nil {
		utils.InternalServerErr(w, err)
	}
	logger, err := middleware.GetLogger(r.Context())
	if err != nil {
		utils.InternalServerErr(w, err)
	}
	reqId, err := middleware.GetReqId(r.Context())
	if err != nil {
		utils.InternalServerErr(w, err)
	}

	logger.Log.Infow("GET_ALL_EXERCISES_CALLED", zap.String("REQ_ID", reqId))

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*3)
	defer cancel()

	in := exerpb.GetAllExercisesREq{UserId: claims.UserId}

	res, err := h.exerClient.GetAllExercises(ctx, &in)
	if err != nil {
		myerrors.ErrMatcher(w, err)
		return
	}

	var resp user.AllExercisesResp

	for _, eachExer := range res.AllExericses {
		exer := user.Exercise{
			Id:        eachExer.Id,
			Name:      eachExer.Name,
			BodyPart:  eachExer.BodyPart,
			Equipment: eachExer.Equipment,
			CreatedAt: eachExer.CreatedAt.AsTime(),
			UpdatedAt: eachExer.UpdatedAt.AsTime(),
		}

		resp.Exercises = append(resp.Exercises, exer)
	}

	resp.NumberOfExercises = int(res.NumberOfExercises)

	utils.OkRespWriter(w, resp)
	// logger.Info("USER_REQUESTED_ALL_EXERCISES", "user_id", 1)
	logger.Log.Infow("GET_ALL_EXERCISES_SUCCESSFUL", zap.String("REQ_ID", reqId))
}

func (h *Handler) DeleteExecise(w http.ResponseWriter, r *http.Request) {

	claims, err := middleware.GetClaims(r.Context())
	if err != nil {
		utils.InternalServerErr(w, err)
	}
	logger, err := middleware.GetLogger(r.Context())
	if err != nil {
		utils.InternalServerErr(w, err)
	}
	reqId, err := middleware.GetReqId(r.Context())
	if err != nil {
		utils.InternalServerErr(w, err)
	}

	logger.Log.Infow("DELETE_EXERCISE_CALLED", zap.String("REQ_ID", reqId))

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*3)
	defer cancel()

	var userInput user.ExerciseName

	json.NewDecoder(r.Body).Decode(&userInput)

	validationErr, errOccured := userInput.Validate()
	if errOccured {
		utils.ValidationErrWriter(w, *validationErr)
		return
	}

	in := exerpb.SendExerciseName{
		ExerciseName: userInput.Name,
		UserId:       claims.UserId,
	}

	resp, err := h.exerClient.DeleteExercise(ctx, &in)
	if err != nil {
		myerrors.ErrMatcher(w, err)
		return
	}

	utils.DeletedNotFoundWriter(w, resp)
	logger.Log.Infow("DELETE_EXERCISE_SUCCESSFULL", zap.String("REQ_ID", reqId))
}

func (h *Handler) UpdateExercise(w http.ResponseWriter, r *http.Request) {}
