package models

import "time"

type Exercise struct {
	Id        int       `json:"id,omitempty"`
	Name      string    `json:"name"`
	RestTime  int       `json:"rest_time_in_seconds"`
	BodyPart  string    `json:"body_part"`
	Equipment string    `json:"equipment"`
	CreatedAt time.Time `json:"created_at"`
}

type DelExerciseResp struct {
	Message string `json:"message"`
}

type CreateExerciseResp struct {
	Message string `json:"message"`
	Exercise Exercise `json:"exercise"`
}