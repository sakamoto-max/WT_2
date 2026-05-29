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

type Client struct {
	Conn *grpc.ClientConn
}

func NewEmptyClient() *Client {
	return &Client{}
}

func (c *Client) OpenConnection(targetServiceAddr string) *Client {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, err := grpc.NewClient(targetServiceAddr, opts...)
	if err != nil {
		log.Fatalf("failed to create the client : %v", err)
	}

	return &Client{Conn: conn}
}

func (c *Client) CreatePlanClient() PlanClientIFace {
	planClient := planpb.NewPlanServiceClient(c.Conn)
	return planClient
}
func (c *Client) CreateExerciseClient() ExerClientIface {
	exerClient := exerpb.NewExerciseServiceClient(c.Conn)
	return exerClient
}

// func (c *client) Close() {
// 	c.ConnForExer.Close()
// 	c.ConnForPlan.Close()
// }

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
