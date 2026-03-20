package handlers

import (
	authpb "workout-tracker/proto/shared/auth"
	exerpb "workout-tracker/proto/shared/exercise"
	planpb "workout-tracker/proto/shared/plan"
	trackpb "workout-tracker/proto/shared/tracker"
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



