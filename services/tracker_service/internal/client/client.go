package client

import (
	"log"
	planpb "workout-tracker/proto/shared/plan"
	exerpb "workout-tracker/proto/shared/exercise"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type client struct{
	connForPlan *grpc.ClientConn
	connForExer *grpc.ClientConn
	PlanClient planpb.PlanServiceClient
	ExerClient exerpb.ExerciseServiceClient
}

// type PlanClient struct{}

func New() *client {

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	connForPlan, err := grpc.NewClient("localhost:6002", opts...)
	if err != nil {
		log.Fatalf("failed to create the client : %v", err)
	}

	planClient := planpb.NewPlanServiceClient(connForPlan)

	connForExer, err := grpc.NewClient("localhost:6003", opts...)
	if err != nil {
		log.Fatalf("failed to create the client : %v", err)
	}
	
	exerClient := exerpb.NewExerciseServiceClient(connForExer)

	return &client{
		connForPlan: connForPlan,
		connForExer: connForExer,
		PlanClient: planClient,
		ExerClient: exerClient,
	}
}


func (c *client) Close() {
	c.connForExer.Close()
	c.connForPlan.Close()
}


