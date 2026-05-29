package client

import (
	"context"
	"log"
	exerpb "github.com/sakamoto-max/wt_2_proto/shared/exercise"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

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

func (c *Client) CreateExerciseClient() ExerClientIface {
	exerClient := exerpb.NewExerciseServiceClient(c.Conn)
	return exerClient
}
type ExerClientIface interface {
	GetExerciseName(ctx context.Context, in *exerpb.SendExerciseID, opts ...grpc.CallOption) (*exerpb.GetExerciseNameResp, error)
	ExerciseExistsReturnId(ctx context.Context, in *exerpb.SendExerciseName, opts ...grpc.CallOption) (*exerpb.ExerciseExistsReturnIdResp, error)
}

// used - getExerciseName  ExerciseExistsReturnId


// CreateExercise(ctx context.Context, in *exerpb.CreateExerciseReq, opts ...grpc.CallOption) (*exerpb.CreateExerciseResp, error)
// DeleteExercise(ctx context.Context, in *exerpb.SendExerciseName, opts ...grpc.CallOption) (*exerpb.DeleteExerciseResp, error)
// ExerciseExistsInMainReturnEveryThing(ctx context.Context, in *exerpb.SendExerciseName, opts ...grpc.CallOption) (*exerpb.SendEverythingResp, error)
// GetAllExercises(ctx context.Context, in *exerpb.GetAllExercisesREq, opts ...grpc.CallOption) (*exerpb.GetAllExercisesResp, error)
// GetHealth(ctx context.Context, in *exerpb.GetHealthReq, opts ...grpc.CallOption) (*exerpb.GetHealthResp, error)
// GetOneExercise(ctx context.Context, in *exerpb.SendExerciseName, opts ...grpc.CallOption) (*exerpb.OneExerciseResp, error)
// PING(ctx context.Context, in *exerpb.PingExerReq, opts ...grpc.CallOption) (*exerpb.PingExerResp, error)
// ExercisesExistsReturnsIds(ctx context.Context, in *exerpb.SendExerciseNames, opts ...grpc.CallOption) (*exerpb.ExerciseIds, error)