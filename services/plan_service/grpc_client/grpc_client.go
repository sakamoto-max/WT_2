package grpcclient

import (
	"log"
	exerpb "workout-tracker/proto/shared/exercise"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ExerciseClient struct{
	conn *grpc.ClientConn
	Client exerpb.ExerciseServiceClient
}

func NewExerciseServiceClient() *ExerciseClient {
	return &ExerciseClient{}
}

func New() *ExerciseClient {
	var opts []grpc.DialOption

	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, err := grpc.NewClient("localhost:6003", opts...)
	if err != nil {
		log.Fatalf("failed to create the client : %v", err)
	}

	client := exerpb.NewExerciseServiceClient(conn)

	return &ExerciseClient{
		conn: conn,
		Client: client,
	}
}

func (e *ExerciseClient) Close() {
	e.conn.Close()
}