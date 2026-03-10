package client

import (
	"context"
	exerpb "workout-tracker/proto/shared/exercise"

	"google.golang.org/grpc"
)

// import (
// 	"context"
// 	"fmt"
// 	"log"
// 	"time"

// 	pb "plan_service/protobuff"

// 	"google.golang.org/grpc"
// 	"google.golang.org/grpc/credentials/insecure"
// )

// func main() {

// 	var opts []grpc.DialOption

// 	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

// 	conn, err := grpc.NewClient("localhost:5008", opts...)
// 	if err != nil {
// 		log.Fatalf("failed to create the client : %v", err)
// 	}

// 	defer conn.Close()

// 	client := pb.NewExerciseServiceClient(conn)

// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
// 	defer cancel()

// 	req := pb.SendExerciseName{
// 		ExerciseName: "Push Ups",
// 	}

// 	resp, err := client.ExerciseExistsReturnId(ctx, &req)
// 	if err != nil {
// 		log.Fatalf("error getting resp from the exercise server : %v", err)
// 	}

// 	fmt.Println(resp)
// }

type ExerciseClient struct {
	client exerpb.ExerciseServiceClient
}

func NewExerciseServiceClient(conn *grpc.ClientConn) *ExerciseClient {
	return &ExerciseClient{
		client: exerpb.NewExerciseServiceClient(conn),
	}
}

func (e *ExerciseClient) ExerciseExistsReturnId(ctx context.Context, req *exerpb.SendExerciseName) (*exerpb.ExerciseExistsReturnIdResp, error) {
	return e.client.ExerciseExistsReturnId(ctx, req)
}

