package grpcclient

import (
	// "context"
	// "context"
	"log"
	"os"
	// "time"

	// "time"
	authpb "workout-tracker/proto/shared/auth"
	exerpb "workout-tracker/proto/shared/exercise"
	planpb "workout-tracker/proto/shared/plan"
	trackpb "workout-tracker/proto/shared/tracker"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type grpcClient struct {
	AuthClient  authpb.AuthServiceClient
	PlanClient  planpb.PlanServiceClient
	ExerClient  exerpb.ExerciseServiceClient
	TrackClient trackpb.TrackerServiceClient
}

func NewgrpcClient() *grpcClient {
	return &grpcClient{}
}

func (g *grpcClient) ConnectToClients() *grpcClient {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	connForAuth, err := grpc.NewClient(os.Getenv("AUTH_GRPC_CLIENT_ADDR"), opts...)
	if err != nil {
		log.Fatalf("error creating auth client : %v", err)
	}

	connForPlan, err := grpc.NewClient(os.Getenv("PLAN_GRPC_CLIENT_ADDR"), opts...)
	if err != nil {
		log.Fatalf("error creating plan client : %v", err)
	}

	connForExer, err := grpc.NewClient(os.Getenv("EXER_GRPC_CLIENT_ADDR"), opts...)
	if err != nil {
		log.Fatalf("error creating exer client : %v", err)
	}

	connForTracker, err := grpc.NewClient(os.Getenv("TRACKER_GRPC_CLIENT_ADDR"), opts...)
	if err != nil {
		log.Fatalf("error creating tracker client : %v", err)
	}

	authClient := authpb.NewAuthServiceClient(connForAuth)
	planClient := planpb.NewPlanServiceClient(connForPlan)
	ExerClient := exerpb.NewExerciseServiceClient(connForExer)
	TrackerClient := trackpb.NewTrackerServiceClient(connForTracker)

	return &grpcClient{
		AuthClient:  authClient,
		PlanClient:  planClient,
		ExerClient:  ExerClient,
		TrackClient: TrackerClient,
	}

}

// func (g *grpcClient) PingAll() {

// 	// ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
// 	// defer cancel()

// 	// inOne := authpb.PINGreq{}
// 	// _, err := g.AuthClient.PING(ctx, &inOne)
// 	// if err != nil {
// 	// 	log.Fatalf("auth service is not up")
// 	// }

// 	// inTwo := planpb.PingPlanReq{}
// 	// _, err = g.PlanClient.PING(ctx, &inTwo)
// 	// if err != nil {
// 	// 	log.Fatalf("plan service is not up")

// 	// }

// 	// inThree := exerpb.PingExerReq{}
// 	// _, err = g.ExerClient.PING(ctx, &inThree)
// 	// if err != nil {
// 	// 	log.Fatalf("exercise service is not up")

// 	// }

// 	// inFour := trackpb.PingTrackReq{}
// 	// _, err = g.TrackClient.PING(ctx, &inFour)
// 	// if err != nil {
// 	// 	log.Fatalf("tracker service is not up")
// 	// }
// }
