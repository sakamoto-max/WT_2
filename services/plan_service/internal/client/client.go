package client

import (
	"context"
	"fmt"
	"log"

	"github.com/sakamoto-max/wt_2_pkg/logger"
	exerpb "github.com/sakamoto-max/wt_2_proto/shared/exercise"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	Conn *grpc.ClientConn
}

func NewEmptyClient() *Client {
	return &Client{}
}

func (c *Client) OpenConnection(targetServiceAddr string) *grpc.ClientConn {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, err := grpc.NewClient(targetServiceAddr, opts...)
	if err != nil {
		log.Fatalf("failed to create the client : %v", err)
	}

	return conn
}

func NewConn(targetServerHostName string, targetServerAddr string, logger *logger.MyLogger) *grpc.ClientConn {
	// make conn
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	url := fmt.Sprintf("%s:%s", targetServerHostName, targetServerAddr)

	conn, err := grpc.NewClient(url, opts...)
	if err != nil {
		logger.Log.Fatalf("failed to connect to the client", zap.Error(err))
	}

	return conn
}

func CreateExerciseClient(conn *grpc.ClientConn) ExerClientIface {
	exerClient := exerpb.NewExerciseServiceClient(conn)
	return exerClient
}

type ExerClientIface interface {
	GetExerciseName(ctx context.Context, in *exerpb.SendExerciseID, opts ...grpc.CallOption) (*exerpb.GetExerciseNameResp, error)
	ExerciseExistsReturnId(ctx context.Context, in *exerpb.SendExerciseName, opts ...grpc.CallOption) (*exerpb.ExerciseExistsReturnIdResp, error)
}
