package domain

import (
	"time"

	planpb "github.com/sakamoto-max/wt_2_proto/shared/plan"
)

type Plan struct {
	UserId        string    `json:"userId,omitempty"`
	PlanId        string    `json:"planId,omitempty"`
	PlanName      string    `json:"plan_name,omitempty"`
	ExerciseNames *[]string `json:"exerciseNames,omitempty"`
	ExerciseIds   *[]string `json:"exerciseIds,omitempty"`
	CreatedAt     time.Time `json:"createdAt,omitempty"`
	UpdatedAt     time.Time `json:"updatedAt,omitempty"`
}

type CreatePlan struct {
	UserId        string
	PlanName      string
	ExerciseNames *[]string
	ExerciseIds   *[]string
}

func ToCreatePlan(in *planpb.CreatePlanReq) CreatePlan {
	return CreatePlan{
		UserId:        in.UserId,
		PlanName:      in.PlanName,
		ExerciseNames: &in.ExerciseNames,
	}
}

type GetPlan struct {
	UserId   string
	PlanName string
}

// func ToGetPlan(in *planpb.GetPlanByNameReq) GetPlan {
// 	return GetPlan{
// 		UserId:   in.UserId,
// 		PlanName: in.PlanName,
// 	}
// }

type PlanReq interface {
	GetPlanName() string
	GetUserId() string
}

func ToGetPlan[T PlanReq](userSentData T) GetPlan {
	planName := userSentData.GetPlanName()
	UserId := userSentData.GetUserId()

	return GetPlan{
		UserId:   UserId,
		PlanName: planName,
	}
}
