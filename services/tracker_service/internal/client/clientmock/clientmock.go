package clientmock

import (
	"context"
	"fmt"

	exerpb "github.com/sakamoto-max/wt_2_proto/shared/exercise"
	planpb "github.com/sakamoto-max/wt_2_proto/shared/plan"
	"google.golang.org/grpc"
)

type PlanClientMock struct {
	Down      bool
	PlanExits bool
}

func (p *PlanClientMock) GetEmptyPlanId(ctx context.Context, in *planpb.SendUserID, opts ...grpc.CallOption) (*planpb.EmptyPlanIdResp, error) {
	if p.Down {
		return nil, fmt.Errorf("client is down")
	}

	if !p.PlanExits {
		return nil, fmt.Errorf("plan doesn't exist")
	}

	return &planpb.EmptyPlanIdResp{EmptyPlanId: "123"}, nil
}
func (p *PlanClientMock) GetPlanByName(ctx context.Context, in *planpb.GetPlanByNameReq, opts ...grpc.CallOption) (*planpb.GetPlanByNameResp, error) {
	if p.Down {
		return nil, fmt.Errorf("client is down")
	}

	if !p.PlanExits {
		return nil, fmt.Errorf("plan doesn't exist")
	}

	return &planpb.GetPlanByNameResp{PlanId: "123", ExerciseNames: []string{"exer_1", "exer_2", "exer_3"}}, nil
}

type ExerClientMock struct{
	Down bool
	ExerciseExists bool
}

func (e *ExerClientMock) ExerciseExistsReturnId(ctx context.Context, in *exerpb.SendExerciseName, opts ...grpc.CallOption) (*exerpb.ExerciseExistsReturnIdResp, error) {
	fmt.Println("client started")
	if e.Down {
		fmt.Println("returnign err")
		return nil, fmt.Errorf("client is down")
	}
	
	if !e.ExerciseExists {
		return nil, fmt.Errorf("exercise doesn't exist")
	}
	
	fmt.Println("client ended")
	return &exerpb.ExerciseExistsReturnIdResp{ExerciseId: "123"}, nil
}
