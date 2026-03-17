package handlers

import (
	"api_gateway/responses"
	"api_gateway/user"
	"context"
	"encoding/json"
	"net/http"
	"time"
	myerrors "wt/pkg/my_errors"
	"wt/pkg/utils"
	authpb "workout-tracker/proto/shared/auth"
	exerpb "workout-tracker/proto/shared/exercise"
	planpb "workout-tracker/proto/shared/plan"
	trackpb "workout-tracker/proto/shared/tracker"
	token "wt/pkg/shared"
)

type Handler struct {
	authClient  authpb.AuthServiceClient
	planClient  planpb.PlanServiceClient
	exerClient  exerpb.ExerciseServiceClient
	trackClient trackpb.TrackerServiceClient
}

func NewHandler(authClient authpb.AuthServiceClient, planClient planpb.PlanServiceClient, exerClient exerpb.ExerciseServiceClient, trackClient trackpb.TrackerServiceClient) *Handler {
	return &Handler{authClient: authClient, planClient: planClient, exerClient: exerClient, trackClient: trackClient}
}


// ERRS :
// 1. user already exists - name, email -> bad req
// 2. auth server isn't responding -> internal server -> show the user 500
// 3. plan server isn't responding -> internal server -> show the user 500
// 4. token err
//       - token is exp
//       - token is missing
// 5. auth is responding but plan isn't working

func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*3)
	defer cancel()

	// parse the data from the user
	var userInput user.Signup

	json.NewDecoder(r.Body).Decode(&userInput)

	// validate it
	validationErrs, errOccured := userInput.Validate()
	if errOccured {
		utils.ValidationErrWriter(w, validationErrs)
		return
	}

	// make a new req
	in := authpb.UserSignUpReq{}
	in.Email = userInput.Email
	in.Name = userInput.Name
	in.Password = userInput.Password

	resp, err := h.authClient.UserSignUp(ctx, &in)
	if err != nil {
		myerrors.ErrMatcher(w, err)
		return
	}

	re := responses.SignUpResp{
		Name:      resp.Name,
		Email:     resp.Email,
		Role:      resp.Role,
		CreatedAt: resp.CreatedAt.AsTime(),
	}

	utils.CreatedWriter(w, re)
}
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	var userInput user.Login

	json.NewDecoder(r.Body).Decode(&userInput)

	validationErr, errOccured := userInput.Validate()
	if errOccured {
		utils.ValidationErrWriter(w, validationErr)
		return
	}

	in := authpb.UserLoginReq{
		Email:    userInput.Email,
		Password: userInput.Password,
	}

	resp, err := h.authClient.UserLogin(ctx, &in)
	if err != nil {

	}

	utils.OkRespWriter(w, resp)
}
func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	t := token.JwtToken{}
	claims, ok := t.GetClaimsFromContext(r.Context())
	if !ok {
		utils.InternalServerErr(w, myerrors.ErrGettingClaims)
	}

	in := authpb.SendUserId{
		UserId: int64(claims.UserId),
	}

	resp, err := h.authClient.UserLogOut(ctx, &in)
	if err != nil {
		utils.BadReq(w, err)
		return
	}

	utils.OkRespWriter(w, resp)
}
func (h *Handler) GetNewAccessToken(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*3)
	defer cancel()

	var userInput user.UUIDReader

	json.NewDecoder(r.Body).Decode(&userInput)

	validationErr, errOccured := userInput.Validate()
	if errOccured {
		utils.ValidationErrWriter(w, validationErr)
		return
	}

	in := authpb.SendUUID{
		UUID: userInput.UUID,
	}

	resp, err := h.authClient.GetNewAccessToken(ctx, &in)
	if err != nil {
		utils.BadReq(w, err)
		return
	}

	utils.CreatedWriter(w, resp)
}
func (h *Handler) GetAllExercises(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*3)
	defer cancel()

	in := exerpb.GetAllExercisesREq{}

	resp, err := h.exerClient.GetAllExercises(ctx, &in)
	if err != nil {
		utils.BadReq(w, err)
	}

	utils.OkRespWriter(w, resp)
}
func (h *Handler) GetExerciseByName(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*3)
	defer cancel()

	var userInput user.ExerciseName

	json.NewDecoder(r.Body).Decode(&userInput)

	validationErr, errOccured := userInput.Validate()
	if errOccured {
		utils.ValidationErrWriter(w, validationErr)
		return
	}

	in := exerpb.SendExerciseName{
		ExerciseName: userInput.Name,
	}

	resp, err := h.exerClient.GetOneExercise(ctx, &in)
	if err != nil {
		utils.BadReq(w, err)
		return
	}

	utils.OkRespWriter(w, resp)
}
func (h *Handler) CreateExercise(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*3)
	defer cancel()

	var userInput user.Exercise

	json.NewDecoder(r.Body).Decode(&userInput)

	validationErr, errOccured := userInput.Validate()
	if errOccured {
		utils.ValidationErrWriter(w, validationErr)
		return
	}

	in := exerpb.CreateExerciseReq{
		ExerciseName: userInput.Name,
		BodyPart:     userInput.BodyPart,
		Equipment:    userInput.Equipment,
		RestTime:     int64(userInput.RestTime),
	}

	resp, err := h.exerClient.CreateExercise(ctx, &in)
	if err != nil {
		utils.BadReq(w, err)
		return
	}

	utils.CreatedWriter(w, resp)
}
func (h *Handler) DeleteExecise(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*3)
	defer cancel()

	var userInput user.ExerciseName

	json.NewDecoder(r.Body).Decode(&userInput)

	validationErr, errOccured := userInput.Validate()
	if errOccured {
		utils.ValidationErrWriter(w, validationErr)
		return
	}

	in := exerpb.SendExerciseName{
		ExerciseName: userInput.Name,
	}

	resp, err := h.exerClient.DeleteExercise(ctx, &in)
	if err != nil {
		utils.BadReq(w, err)
		return
	}

	utils.DeletedNotFoundWriter(w, resp)
}
func (h *Handler) CheckHealthPlan(w http.ResponseWriter, r *http.Request) {
	// ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	// defer cancel()

	// t := token.JwtToken{}
	// claims, ok := t.GetClaimsFromContext(r.Context())
	// if !ok {
	// 	utils.InternalServerErr(w, myerrors.ErrGettingClaims)
	// }

}
func (h *Handler) CreatePlan(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	t := token.JwtToken{}
	claims, ok := t.GetClaimsFromContext(r.Context())
	if !ok {
		utils.InternalServerErr(w, myerrors.ErrGettingClaims)
	}

	var userInput user.Plan

	json.NewDecoder(r.Body).Decode(&userInput)

	validationErrs, errOccured := userInput.Validate()
	if errOccured {
		utils.ValidationErrWriter(w, validationErrs)
		return
	}

	in := planpb.CreatePlanReq{
		UserId:        int64(claims.UserId),
		PlanName:      userInput.PlanName,
		ExerciseNames: userInput.Exercises,
	}

	resp, err := h.planClient.CreatePlan(ctx, &in)
	if err != nil {
		utils.BadReq(w, err)
		return
	}

	utils.CreatedWriter(w, resp)
}
func (h *Handler) GetAllPlans(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	t := token.JwtToken{}
	claims, ok := t.GetClaimsFromContext(r.Context())
	if !ok {
		utils.InternalServerErr(w, myerrors.ErrGettingClaims)
	}

	in := planpb.GetAllPlansReq{
		UserId: int64(claims.UserId),
	}

	resp, err := h.planClient.GetAllPlans(ctx, &in)
	if err != nil {
		utils.BadReq(w, err)
		return
	}

	utils.OkRespWriter(w, resp)
}
func (h *Handler) GetPLanByName(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	t := token.JwtToken{}
	claims, ok := t.GetClaimsFromContext(r.Context())
	if !ok {
		utils.InternalServerErr(w, myerrors.ErrGettingClaims)
	}

	var userInput user.PlanName

	json.NewDecoder(r.Body).Decode(&userInput)

	validationErrs, errOccured := userInput.Validate()
	if errOccured {
		utils.ValidationErrWriter(w, validationErrs)
		return
	}

	in := planpb.GetPlanByNameReq{
		UserId:   int64(claims.UserId),
		PlanName: userInput.PlanName,
	}

	resp, err := h.planClient.GetPlanByName(ctx, &in)
	if err != nil {
		utils.BadReq(w, err)
		return
	}

	utils.OkRespWriter(w, resp)
}
func (h *Handler) AddExercisesToPlan(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	t := token.JwtToken{}
	claims, ok := t.GetClaimsFromContext(r.Context())
	if !ok {
		utils.InternalServerErr(w, myerrors.ErrGettingClaims)
	}

	var userInput user.Plan

	json.NewDecoder(r.Body).Decode(&userInput)

	validationErrs, errOccured := userInput.Validate()
	if errOccured {
		utils.ValidationErrWriter(w, validationErrs)
		return
	}

	in := planpb.PlanReq{
		UserId:        int64(claims.UserId),
		PlanName:      userInput.PlanName,
		ExerciseNames: userInput.Exercises,
	}

	resp, err := h.planClient.AddExercisesToPlan(ctx, &in)
	if err != nil {
		utils.BadReq(w, err)
		return
	}

	utils.OkRespWriter(w, resp)
}
func (h *Handler) DeleteExerciseFromPlan(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	t := token.JwtToken{}
	claims, ok := t.GetClaimsFromContext(r.Context())
	if !ok {
		utils.InternalServerErr(w, myerrors.ErrGettingClaims)
	}

	var userInput user.Plan

	json.NewDecoder(r.Body).Decode(&userInput)

	validationErrs, errOccured := userInput.Validate()
	if errOccured {
		utils.ValidationErrWriter(w, validationErrs)
		return
	}

	in := planpb.PlanReq{
		UserId:        int64(claims.UserId),
		PlanName:      userInput.PlanName,
		ExerciseNames: userInput.Exercises,
	}

	resp, err := h.planClient.DeleteExercisesFromPlan(ctx, &in)
	if err != nil {
		utils.BadReq(w, err)
		return
	}

	utils.OkRespWriter(w, resp)
}
func (h *Handler) DeletePlan(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	t := token.JwtToken{}
	claims, ok := t.GetClaimsFromContext(r.Context())
	if !ok {
		utils.InternalServerErr(w, myerrors.ErrGettingClaims)
	}

	var userInput user.PlanName

	json.NewDecoder(r.Body).Decode(&userInput)

	validationErrs, errOccured := userInput.Validate()
	if errOccured {
		utils.ValidationErrWriter(w, validationErrs)
		return
	}

	in := planpb.DeletePlanReq{
		UserId:   int64(claims.UserId),
		PlanName: userInput.PlanName,
	}

	resp, err := h.planClient.DeletePlan(ctx, &in)
	if err != nil {
		utils.BadReq(w, err)
		return
	}

	utils.OkRespWriter(w, resp)
}
func (h *Handler) StartEmptyWorkout(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	t := token.JwtToken{}
	claims, ok := t.GetClaimsFromContext(r.Context())
	if !ok {
		utils.InternalServerErr(w, myerrors.ErrGettingClaims)
	}

	in := trackpb.StartEmptyWorkoutReq{
		UserId: int64(claims.UserId),
	}

	resp, err := h.trackClient.StartEmptyWorkout(ctx, &in)
	if err != nil {
		utils.BadReq(w, err)
		return
	}

	utils.CreatedWriter(w, resp)
}
func (h *Handler) StartWorkoutWithPlan(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	t := token.JwtToken{}
	claims, ok := t.GetClaimsFromContext(r.Context())
	if !ok {
		utils.InternalServerErr(w, myerrors.ErrGettingClaims)
	}

	var userInput user.PlanName

	json.NewDecoder(r.Body).Decode(&userInput)

	validationErrs, errOccured := userInput.Validate()
	if errOccured {
		utils.ValidationErrWriter(w, validationErrs)
		return
	}

	in := trackpb.StartWorkoutWithPlanReq{
		UserId:   int64(claims.UserId),
		PlanName: userInput.PlanName,
	}

	resp, err := h.trackClient.StartWorkoutWithPlan(ctx, &in)
	if err != nil {
		utils.BadReq(w, err)
		return
	}

	utils.CreatedWriter(w, resp)
}
func (h *Handler) EndWorkout(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	t := token.JwtToken{}
	claims, ok := t.GetClaimsFromContext(r.Context())
	if !ok {
		utils.InternalServerErr(w, myerrors.ErrGettingClaims)
	}

	var userInput user.Tracker

	json.NewDecoder(r.Body).Decode(&userInput)

	validationErrs, errOccured := userInput.Validate()
	if errOccured {
		utils.ValidationErrWriter(w, validationErrs)
		return
	}

	in := trackpb.EndWorkoutReq{}
	in.UserId = int64(claims.UserId)

	for _, allExercises := range userInput.Workout {
		allExer := trackpb.TrackerForEachExer{}
		allExer.ExerciseId = int64(allExercises.ExerciseId)

		for _, eachExer := range allExercises.RepsWeight {
			rw := trackpb.SetsAndReps{}
			rw.Reps = int64(eachExer.Reps)
			rw.Weight = int64(eachExer.Weight)

			allExer.SetsAndReps = append(allExer.SetsAndReps, &rw)
		}

		in.AllExerices = append(in.AllExerices, &allExer)
	}

	resp, err := h.trackClient.EndWorkout(ctx, &in)
	if err != nil {
		utils.BadReq(w, err)
		return
	}

	utils.OkRespWriter(w, resp)
}
