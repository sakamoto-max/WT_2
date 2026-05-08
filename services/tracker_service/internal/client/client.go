package client

import (
	"log"
	"os"
	"fmt"
	planpb "github.com/sakamoto-max/wt_2_proto/shared/plan"
	exerpb "github.com/sakamoto-max/wt_2_proto/shared/exercise"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type client struct{
	connForPlan *grpc.ClientConn
	connForExer *grpc.ClientConn
	PlanClient planpb.PlanServiceClient
	ExerClient exerpb.ExerciseServiceClient
}

func New() *client {

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	fmt.Println("the plan is ", os.Getenv("PLAN_GRPC_SERVER_ADDR"))

	connForPlan, err := grpc.NewClient(os.Getenv("PLAN_GRPC_SERVER_ADDR"), opts...)
	if err != nil {
		log.Fatalf("failed to create the client : %v", err)
	}

	planClient := planpb.NewPlanServiceClient(connForPlan)

	fmt.Println("the exer is", os.Getenv("EXER_GRPC_SERVER_ADDR"))

	connForExer, err := grpc.NewClient(os.Getenv("EXER_GRPC_SERVER_ADDR"), opts...)
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


