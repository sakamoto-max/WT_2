package grpcclient

import (
	"fmt"
	"os"
	authpb "github.com/sakamoto-max/wt_2-proto/shared/auth"
	planpb "github.com/sakamoto-max/wt_2-proto/shared/plan"
	exerpb "github.com/sakamoto-max/wt_2-proto/shared/exercise"
	trackpb "github.com/sakamoto-max/wt_2-proto/shared/tracker"
	"github.com/sakamoto-max/wt_2-pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type grpcClient struct {
	authConn  *grpc.ClientConn
	planConn  *grpc.ClientConn
	exerConn  *grpc.ClientConn
	trackConn *grpc.ClientConn

	AuthClient  authpb.AuthServiceClient
	PlanClient  planpb.PlanServiceClient
	ExerClient  exerpb.ExerciseServiceClient
	TrackClient trackpb.TrackerServiceClient
}

func NewgrpcClient() *grpcClient {
	return &grpcClient{}
}

func (g *grpcClient) ConnectToClients(logger *logger.MyLogger) *grpcClient {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	connForAuth, err := grpc.NewClient(os.Getenv("AUTH_GRPC_CLIENT_ADDR"), opts...)
	if err != nil {
		logger.Log.Fatalf("error creating auth client : %v", err)
	}

	connForPlan, err := grpc.NewClient(os.Getenv("PLAN_GRPC_CLIENT_ADDR"), opts...)
	if err != nil {
		logger.Log.Fatalf("error creating plan client : %v", err)
	}

	connForExer, err := grpc.NewClient(os.Getenv("EXER_GRPC_CLIENT_ADDR"), opts...)
	if err != nil {
		logger.Log.Fatalf("error creating exer client : %v", err)
	}

	connForTracker, err := grpc.NewClient(os.Getenv("TRACKER_GRPC_CLIENT_ADDR"), opts...)
	if err != nil {
		logger.Log.Fatalf("error creating tracker client : %v", err)
	}

	authClient := authpb.NewAuthServiceClient(connForAuth)
	planClient := planpb.NewPlanServiceClient(connForPlan)
	ExerClient := exerpb.NewExerciseServiceClient(connForExer)
	TrackerClient := trackpb.NewTrackerServiceClient(connForTracker)

	return &grpcClient{
		authConn:  connForAuth,
		planConn:  connForPlan,
		exerConn:  connForExer,
		trackConn: connForTracker,

		AuthClient:  authClient,
		PlanClient:  planClient,
		ExerClient:  ExerClient,
		TrackClient: TrackerClient,
	}
}

func (g *grpcClient) Close() error {
	if err := g.authConn.Close(); err != nil {
		return fmt.Errorf("error occured while closing auth client : %w", err)
	}
	if err := g.planConn.Close(); err != nil {
		return fmt.Errorf("error occured while closing auth client : %w", err)
	}
	if err := g.trackConn.Close(); err != nil {
		return fmt.Errorf("error occured while closing auth client : %w", err)
	}
	if err := g.exerConn.Close(); err != nil {
		return fmt.Errorf("error occured while closing auth client : %w", err)
	}

	return nil
}
