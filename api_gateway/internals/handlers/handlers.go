package handlers

import (
	"context"
	"fmt"
	"net/http"
	"time"
	"github.com/sakamoto-max/wt_2/api_gateway/internals/domain"
	"github.com/sakamoto-max/wt_2/api_gateway/internals/utils"
	authpb "github.com/sakamoto-max/wt_2_proto/shared/auth"
	exerpb "github.com/sakamoto-max/wt_2_proto/shared/exercise"
	planpb "github.com/sakamoto-max/wt_2_proto/shared/plan"
	trackpb "github.com/sakamoto-max/wt_2_proto/shared/tracker"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	PlanService     string = "plan_service"
	AuthService     string = "auth_service"
	TrackerService  string = "tracker_service"
	ExerciseService string = "exercise_service"
	EmailService    string = "email_service"

	Healthy     string = "healthy"
	UnHealthy   string = "unhealthy"
	Unavailable string = "unavailable"
)

type Handler struct {
	authClient  authpb.AuthServiceClient
	planClient  planpb.PlanServiceClient
	exerClient  exerpb.ExerciseServiceClient
	trackClient trackpb.TrackerServiceClient
}

func NewHandler(authClient authpb.AuthServiceClient, planClient planpb.PlanServiceClient, exerClient exerpb.ExerciseServiceClient, trackClient trackpb.TrackerServiceClient) *Handler {
	return &Handler{authClient: authClient, planClient: planClient, exerClient: exerClient, trackClient: trackClient}
}

func (h *Handler) GetHealth(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	var health domain.OverAllStatus

	authStatus := domain.OneServiceStatus{}
	authStatus.ServiceName = AuthService

	inForAuth := authpb.GetHealthReq{}
	auth, err := h.authClient.GetHealth(ctx, &inForAuth)
	if err != nil {
		st, _ := status.FromError(err)
		if st.Code() == codes.Unavailable {
			authStatus.Status = UnHealthy
			authStatus.PGRespTime = nil
			authStatus.RedisRespTime = nil
			fmt.Println("auth status : ", authStatus)
			health.AllServices = append(health.AllServices, authStatus)
		}
	} else {
		if auth.PostgresRespTime == nil || auth.RedisRespTime == nil {
			authStatus.Status = UnHealthy
		} else {
			authStatus.Status = Healthy
		}
		pgGoTime := auth.PostgresRespTime.AsDuration().Seconds()
		redisGoTime := auth.RedisRespTime.AsDuration().Seconds()
		authStatus.PGRespTime = &pgGoTime
		authStatus.RedisRespTime = &redisGoTime
		fmt.Println("auth status : ", authStatus)
		health.AllServices = append(health.AllServices, authStatus)
	}
	///////////////////////////////////////////////////////////////////////////////////////////////
	exerStatus := domain.OneServiceStatus{}
	exerStatus.ServiceName = ExerciseService

	inForExer := exerpb.GetHealthReq{}
	exer, err := h.exerClient.GetHealth(ctx, &inForExer)
	if err != nil {
		st, _ := status.FromError(err)
		if st.Code() == codes.Unavailable {
			exerStatus.Status = UnHealthy
			exerStatus.PGRespTime = nil
			exerStatus.RedisRespTime = nil
			fmt.Println("exer status : ", exerStatus)
			health.AllServices = append(health.AllServices, exerStatus)
		}
	} else {
		if exer.PostgresRespTime == nil || exer.RedisRespTime == nil {
			exerStatus.Status = UnHealthy
		} else {
			exerStatus.Status = Healthy
		}
		pgGoTime := exer.PostgresRespTime.AsDuration().Seconds()
		redisGoTime := exer.RedisRespTime.AsDuration().Seconds()
		exerStatus.PGRespTime = &pgGoTime
		exerStatus.RedisRespTime = &redisGoTime
		fmt.Println("exer status : ", exerStatus)
		health.AllServices = append(health.AllServices, exerStatus)
	}
	///////////////////////////////////////////////////////////////////////////////////////////////
	planStatus := domain.OneServiceStatus{}
	planStatus.ServiceName = PlanService

	inForPlan := planpb.GetHealthReq{}
	plan, err := h.planClient.GetHealth(ctx, &inForPlan)
	if err != nil {
		st, _ := status.FromError(err)
		if st.Code() == codes.Unavailable {
			planStatus.Status = UnHealthy
			planStatus.PGRespTime = nil
			planStatus.RedisRespTime = nil
			fmt.Println("plan status : ", planStatus)
			health.AllServices = append(health.AllServices, planStatus)
		}
	} else {
		if plan.PostgresRespTime == nil || plan.RedisRespTime == nil {
			planStatus.Status = UnHealthy
		} else {
			planStatus.Status = Healthy
		}
		pgGoTime := plan.PostgresRespTime.AsDuration().Seconds()
		redisGoTime := plan.RedisRespTime.AsDuration().Seconds()
		planStatus.PGRespTime = &pgGoTime
		planStatus.RedisRespTime = &redisGoTime
		fmt.Println("plan status : ", planStatus)
		health.AllServices = append(health.AllServices, planStatus)
	}
	////////////////////////////////////////////////////////////////////////////////////////////////
	trackerStatus := domain.OneServiceStatus{}
	trackerStatus.ServiceName = TrackerService

	inForTracker := trackpb.GetHealthReq{}
	tracker, err := h.trackClient.GetHealth(ctx, &inForTracker)
	if err != nil {
		st, _ := status.FromError(err)
		if st.Code() == codes.Unavailable {
			trackerStatus.Status = UnHealthy
			trackerStatus.PGRespTime = nil
			trackerStatus.RedisRespTime = nil
			fmt.Println("tracker status : ", trackerStatus)
			health.AllServices = append(health.AllServices, trackerStatus)
		}
	} else {
		if tracker.PostgresRespTime == nil || tracker.RedisRespTime == nil {
			trackerStatus.Status = UnHealthy
		} else {
			trackerStatus.Status = Healthy
		}

		pgGoTime := tracker.PostgresRespTime.AsDuration().Seconds()
		redisGoTime := tracker.RedisRespTime.AsDuration().Seconds()
		trackerStatus.PGRespTime = &pgGoTime
		trackerStatus.RedisRespTime = &redisGoTime
		fmt.Println("tracker status : ", trackerStatus)
		health.AllServices = append(health.AllServices, trackerStatus)
	}
	//////////////////////////////////////////////////////////////////////////////////////////////

	if authStatus.Status == UnHealthy || exerStatus.Status == UnHealthy || trackerStatus.Status == UnHealthy || planStatus.Status == UnHealthy {
		health.Status = UnHealthy
	} else {
		health.Status = Healthy
	}

	utils.OkRespWriter(w, health)

	// {
	// 	"status" : "healthy"/"unhealthy",
	// 	"all_services" : [
	// 		{
	// 			"service_name" : "auth_service"
	// 			"status" : "healthy"/"unhealthy",
	// 			"postgres_response_time" : 5.98039,
	// 			"redis_response_time" : 2.3023
	// 		},
	// 		{
	// 			"service_name" : "plan_service"
	// 			"status" : "healthy"/"unhealthy",
	// 			"postgres_response_time" : 5.98039,
	// 			"redis_response_time" : 2.3023
	// 		},
	// 		{
	// 			"service_name" : "tracker_service"
	// 			"status" : "healthy"/"unhealthy",
	// 			"postgres_response_time" : 5.98039,
	// 			"redis_response_time" : 2.3023
	// 		},
	// 		{
	// 			"service_name" : "exer_service"
	// 			"status" : "healthy"/"unhealthy",
	// 			"postgres_response_time" : 5.98039,
	// 			"redis_response_time" : 2.3023
	// 		},
	// 	]
	// }

}
