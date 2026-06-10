package client

import (
	"context"
	// "fmt"
	"log"
	// "os"

	exerpb "github.com/sakamoto-max/wt_2_proto/shared/exercise"
	planpb "github.com/sakamoto-max/wt_2_proto/shared/plan"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type client struct {
	ConnForPlan *grpc.ClientConn
	ConnForExer *grpc.ClientConn
	PlanClient  planpb.PlanServiceClient
	ExerClient  exerpb.ExerciseServiceClient
}

func OpenConnection(targetServiceAddr string) *grpc.ClientConn {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, err := grpc.NewClient(targetServiceAddr, opts...)
	if err != nil {
		log.Fatalf("failed to create the client : %v", err)
	}

	return conn
}

func CreatePlanClient(conn *grpc.ClientConn) PlanClientIFace {
	planClient := planpb.NewPlanServiceClient(conn)
	return planClient
}
func CreateExerciseClient(conn *grpc.ClientConn) ExerClientIface {
	exerClient := exerpb.NewExerciseServiceClient(conn)
	return exerClient
}

type PlanClientIFace interface {
	GetEmptyPlanId(ctx context.Context, in *planpb.SendUserID, opts ...grpc.CallOption) (*planpb.EmptyPlanIdResp, error)
	GetPlanByName(ctx context.Context, in *planpb.GetPlanByNameReq, opts ...grpc.CallOption) (*planpb.GetPlanByNameResp, error)
}

// GetEmptyPlanId
// GetPlanByName

type ExerClientIface interface {
	ExerciseExistsReturnId(ctx context.Context, in *exerpb.SendExerciseName, opts ...grpc.CallOption) (*exerpb.ExerciseExistsReturnIdResp, error)
}

// exerciseExistsReturnsId

// type planClient struct {
// 	client planpb.PlanServiceClient
// }
