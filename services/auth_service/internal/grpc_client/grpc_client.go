package grpcclient

import (
	// "context"
	"log"
	planpb "workout-tracker/proto/shared/plan"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type PlanClient struct {
	conn   *grpc.ClientConn
	Client planpb.PlanServiceClient
}

func NewPlanClient() *PlanClient {
	return &PlanClient{}
}
func New() *PlanClient {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, err := grpc.NewClient("localhost:6002", opts...)
	if err != nil {
		log.Fatalf("failed to create the client : %v", err)
	}

	client := planpb.NewPlanServiceClient(conn)


	return &PlanClient{
		conn: conn,
		Client: client,
	}
}
func (p *PlanClient) Close() {
	if err := p.conn.Close(); err != nil{
		log.Printf("error closing the client : %v", err)
	}
}
