package grpcclient

import (
	"fmt"
	"github.com/sakamoto-max/wt_2/api_gateway/internals/config"
	authpb "github.com/sakamoto-max/wt_2_proto/shared/auth"
	exerpb "github.com/sakamoto-max/wt_2_proto/shared/exercise"
	planpb "github.com/sakamoto-max/wt_2_proto/shared/plan"
	trackpb "github.com/sakamoto-max/wt_2_proto/shared/tracker"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GrpcClient struct {
	authConn  *grpc.ClientConn
	planConn  *grpc.ClientConn
	exerConn  *grpc.ClientConn
	trackConn *grpc.ClientConn

	AuthClient  authpb.AuthServiceClient
	PlanClient  planpb.PlanServiceClient
	ExerClient  exerpb.ExerciseServiceClient
	TrackClient trackpb.TrackerServiceClient
}

func ConnectToClients(config config.Config) *GrpcClient {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	authUrl := fmt.Sprintf("%v:%v", config.AuthServerConfig.Host, config.AuthServerConfig.Addr)

	connForAuth, err := grpc.NewClient(authUrl)
	if err != nil {
		config.Logger.Log.Fatalf("error creating auth client : %v", err)
	}

	planUrl := fmt.Sprintf("%v:%v", config.PlanServerConfig.Host, config.PlanServerConfig.Addr)

	connForPlan, err := grpc.NewClient(planUrl)
	if err != nil {
		config.Logger.Log.Fatalf("error creating plan client : %v", err)
	}

	exerUrl := fmt.Sprintf("%v:%v", config.ExerServerConfig.Host, config.ExerServerConfig.Addr)

	connForExer, err := grpc.NewClient(exerUrl)
	if err != nil {
		config.Logger.Log.Fatalf("error creating exer client : %v", err)
	}

	trackerUrl := fmt.Sprintf("%v:%v", config.TrackerServerConfig.Host, config.TrackerServerConfig.Addr)

	connForTracker, err := grpc.NewClient(trackerUrl)
	if err != nil {
		config.Logger.Log.Fatalf("error creating tracker client : %v", err)
	}

	authClient := authpb.NewAuthServiceClient(connForAuth)
	planClient := planpb.NewPlanServiceClient(connForPlan)
	ExerClient := exerpb.NewExerciseServiceClient(connForExer)
	TrackerClient := trackpb.NewTrackerServiceClient(connForTracker)

	return &GrpcClient{
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

func (g *GrpcClient) Close() error {
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
