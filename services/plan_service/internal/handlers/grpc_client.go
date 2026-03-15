package handlers

import (
	"context"
	exerpb "workout-tracker/proto/shared/exercise"

	"google.golang.org/grpc"
)


type ExerciseClient struct {
	client exerpb.ExerciseServiceClient
}

func NewExerciseServiceClient(conn *grpc.ClientConn) *ExerciseClient {
	return &ExerciseClient{
		client: exerpb.NewExerciseServiceClient(conn),
	}
}

func (e *ExerciseClient) ExerciseExistsReturnId(ctx context.Context, req *exerpb.SendExerciseName) (*exerpb.ExerciseExistsReturnIdResp, error) {
	return e.client.ExerciseExistsReturnId(ctx, req)
}
