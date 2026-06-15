package mockclient

import (
	"context"
	"fmt"

	exerpb "github.com/sakamoto-max/wt_2_proto/shared/exercise"
	"google.golang.org/grpc"
)

type MockExerClient struct {
	ExerExists bool
	Down       bool
}

func (m *MockExerClient) GetExerciseName(ctx context.Context, in *exerpb.SendExerciseID, opts ...grpc.CallOption) (*exerpb.GetExerciseNameResp, error) {

	if m.Down {
		return nil, fmt.Errorf("service is down")
	}

	if !m.ExerExists {
		return nil, fmt.Errorf("exercise doesn't exist")
	}

	return &exerpb.GetExerciseNameResp{
		ExerciseName: "name",
	}, nil

}

func (m *MockExerClient) ExerciseExistsReturnId(ctx context.Context, in *exerpb.SendExerciseName, opts ...grpc.CallOption) (*exerpb.ExerciseExistsReturnIdResp, error) {
	if m.Down {
		return nil, fmt.Errorf("service is down")
	}

	if !m.ExerExists {
		return nil, fmt.Errorf("exercise doesn't exist")
	}

	return &exerpb.ExerciseExistsReturnIdResp{
		ExerciseId: "123",
	}, nil
}
