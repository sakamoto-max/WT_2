package handlers

import (
	"context"
	"fmt"
	"net/http"
	"time"
	authpb "github.com/sakamoto-max/wt_2-proto/shared/auth"
	exerpb "github.com/sakamoto-max/wt_2-proto/shared/exercise"
	planpb "github.com/sakamoto-max/wt_2-proto/shared/plan"
	trackpb "github.com/sakamoto-max/wt_2-proto/shared/tracker"
	"github.com/sakamoto-max/wt_2-pkg/enum"
	"github.com/sakamoto-max/wt_2/api_gateway/domain"
	"github.com/sakamoto-max/wt_2/api_gateway/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	authStatus.ServiceName = string(enum.AuthService)

	inForAuth := authpb.GetHealthReq{}
	auth, err := h.authClient.GetHealth(ctx, &inForAuth)
	if err != nil {
		st, _ := status.FromError(err)
		if st.Code() == codes.Unavailable {
			authStatus.Status = string(enum.UnHealthy)
			authStatus.PGRespTime = nil
			authStatus.RedisRespTime = nil
			fmt.Println("auth status : ", authStatus)
			health.AllServices = append(health.AllServices, authStatus)
		}
	} else {
		if auth.PostgresRespTime == nil || auth.RedisRespTime == nil {
			authStatus.Status = string(enum.UnHealthy)
		} else {
			authStatus.Status = string(enum.Healthy)
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
	exerStatus.ServiceName = string(enum.ExerciseService)

	inForExer := exerpb.GetHealthReq{}
	exer, err := h.exerClient.GetHealth(ctx, &inForExer)
	if err != nil {
		st, _ := status.FromError(err)
		if st.Code() == codes.Unavailable {
			exerStatus.Status = string(enum.UnHealthy)
			exerStatus.PGRespTime = nil
			exerStatus.RedisRespTime = nil
			fmt.Println("exer status : ", exerStatus)
			health.AllServices = append(health.AllServices, exerStatus)
		}
	} else {
		if exer.PostgresRespTime == nil || exer.RedisRespTime == nil {
			exerStatus.Status = string(enum.UnHealthy)
		} else {
			exerStatus.Status = string(enum.Healthy)
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
	planStatus.ServiceName = string(enum.PlanService)

	inForPlan := planpb.GetHealthReq{}
	plan, err := h.planClient.GetHealth(ctx, &inForPlan)
	if err != nil {
		st, _ := status.FromError(err)
		if st.Code() == codes.Unavailable {
			planStatus.Status = string(enum.UnHealthy)
			planStatus.PGRespTime = nil
			planStatus.RedisRespTime = nil
			fmt.Println("plan status : ", planStatus)
			health.AllServices = append(health.AllServices, planStatus)
		}
	} else {
		if plan.PostgresRespTime == nil || plan.RedisRespTime == nil {
			planStatus.Status = string(enum.UnHealthy)
		} else {
			planStatus.Status = string(enum.Healthy)
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
	trackerStatus.ServiceName = string(enum.TrackerService)

	inForTracker := trackpb.GetHealthReq{}
	tracker, err := h.trackClient.GetHealth(ctx, &inForTracker)
	if err != nil {
		st, _ := status.FromError(err)
		if st.Code() == codes.Unavailable {
			trackerStatus.Status = string(enum.UnHealthy)
			trackerStatus.PGRespTime = nil
			trackerStatus.RedisRespTime = nil
			fmt.Println("tracker status : ", trackerStatus)
			health.AllServices = append(health.AllServices, trackerStatus)
		}
	} else {
		if tracker.PostgresRespTime == nil || tracker.RedisRespTime == nil {
			trackerStatus.Status = string(enum.UnHealthy)
		} else {
			trackerStatus.Status = string(enum.Healthy)
		}

		pgGoTime := tracker.PostgresRespTime.AsDuration().Seconds()
		redisGoTime := tracker.RedisRespTime.AsDuration().Seconds()
		trackerStatus.PGRespTime = &pgGoTime
		trackerStatus.RedisRespTime = &redisGoTime
		fmt.Println("tracker status : ", trackerStatus)
		health.AllServices = append(health.AllServices, trackerStatus)
	}
	//////////////////////////////////////////////////////////////////////////////////////////////
	
	
	
	if authStatus.Status == string(enum.UnHealthy) || exerStatus.Status == string(enum.UnHealthy) || trackerStatus.Status == string(enum.UnHealthy) || planStatus.Status == string(enum.UnHealthy)  {
		health.Status = string(enum.UnHealthy)
	} else {
		health.Status = string(enum.Healthy)
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

// inForAuth := authpb.GetHealthReq{}
// auth, _ := h.authClient.GetHealth(ctx, &inForAuth)

// authStatus := domain.OneServiceStatus{}

// authStatus.ServiceName = string(enum.AuthService)

// if auth.PostgresRespTime == nil || auth.RedisRespTime == nil {
// 	authStatus.Status = string(enum.UnHealthy)
// }else {
// 	authStatus.Status = string(enum.Healthy)
// }

// pgGoTime := auth.PostgresRespTime.AsDuration()
// redisGoTime := auth.RedisRespTime.AsDuration()

// authStatus.PGRespTime = &pgGoTime
// authStatus.RedisRespTime = &redisGoTime

func new() {

}
