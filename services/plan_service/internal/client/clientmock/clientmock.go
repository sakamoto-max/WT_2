package clientmock

import (
	"context"
	"fmt"

	exerpb "github.com/sakamoto-max/wt_2_proto/shared/exercise"
	"google.golang.org/grpc"
)

type ClientMock struct {
	ServerIsDown   bool
	ExerciseExists bool
}

func (c *ClientMock) ExerciseExistsReturnId(ctx context.Context, in *exerpb.SendExerciseName, opts ...grpc.CallOption) (*exerpb.ExerciseExistsReturnIdResp, error) {
	if c.ServerIsDown {
		return nil, fmt.Errorf("server is not responding")
	}

	if !c.ExerciseExists {
		return nil, fmt.Errorf("exercise doesn't exists")
	}

	return &exerpb.ExerciseExistsReturnIdResp{ExerciseId: "123"}, nil
}

func (c *ClientMock) GetExerciseName(ctx context.Context, in *exerpb.SendExerciseID, opts ...grpc.CallOption) (*exerpb.GetExerciseNameResp, error) {
	if c.ServerIsDown {
		return nil, fmt.Errorf("server is not responding")
	}

	if !c.ExerciseExists {
		return nil, fmt.Errorf("exercise doesn't exists")
	}

	return &exerpb.GetExerciseNameResp{ExerciseName: "name"}, nil
}

// not needed