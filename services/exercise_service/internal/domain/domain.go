package domain

import (
	"time"
)

type Exercise struct {
	Id        string    `json:"id" redis:"id"`
	Name      string    `json:"name" redis:"exercise_name"`
	RestTime  int       `json:"restTime" redis:"rest_time"`
	BodyPart  string    `json:"bodyPart" redis:"body_part"`
	Equipment string    `json:"equipment" redis:"equipment"`
	CreatedAt time.Time `json:"createdAt" redis:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" redis:"updated_at"`
}

type DelExerciseResp struct {
	Message string `json:"message"`
}

type CreateExerciseResp struct {
	Message  string   `json:"message"`
	Exercise Exercise `json:"exercise"`
}
