package grpcclient

import (
	"log"
	planpb "workout-tracker/proto/shared/plan"
	exerpb "workout-tracker/proto/shared/exercise"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type PlanClient struct{}

func NewPlanClient() *PlanClient {
	return &PlanClient{}
}


func (p *PlanClient) Connect() planpb.PlanServiceClient {
	
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	
	conn, err := grpc.NewClient("localhost:6002", opts...)
	if err != nil {
		log.Fatalf("failed to create the client : %v", err)
	}
	
	client := planpb.NewPlanServiceClient(conn)
	
	return client
}

type ExerClient struct{}

func NewExerClient() *ExerClient {
	return &ExerClient{}
}

func (e *ExerClient) Connect() exerpb.ExerciseServiceClient {
	
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	
	conn, err := grpc.NewClient("localhost:6003", opts...)
	if err != nil {
		log.Fatalf("failed to create the client : %v", err)
	}
	
	client := exerpb.NewExerciseServiceClient(conn)
	
	return client
}

