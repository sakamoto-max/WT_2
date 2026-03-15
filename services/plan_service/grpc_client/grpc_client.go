package grpcclient

import (
	"log"
	exerpb "workout-tracker/proto/shared/exercise"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ExerciseClient struct {}

func NewExerciseServiceClient() *ExerciseClient {
	return &ExerciseClient{}
}

func (e *ExerciseClient) Connect() exerpb.ExerciseServiceClient {
	var opts []grpc.DialOption

	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, err := grpc.NewClient("localhost:6003", opts...)
	if err != nil {
		log.Fatalf("failed to create the client : %v", err)
	}

	client := exerpb.NewExerciseServiceClient(conn)

	return client
}

