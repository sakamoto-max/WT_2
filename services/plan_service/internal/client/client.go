package client

import (
	"log"
	"os"
	exerpb "github.com/sakamoto-max/wt_2_proto/shared/exercise"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type exerciseClient struct {
	conn   *grpc.ClientConn
	Client exerpb.ExerciseServiceClient
}

func New() *exerciseClient {
	var opts []grpc.DialOption

	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, err := grpc.NewClient(os.Getenv("EXERCISE_GRPC_SERVER_ADDR"), opts...)
	if err != nil {
		log.Fatalf("failed to create the client : %v", err)
	}

	client := exerpb.NewExerciseServiceClient(conn)

	return &exerciseClient{
		conn:   conn,
		Client: client,
	}
}

func (e *exerciseClient) Close() {
	e.conn.Close()
}
