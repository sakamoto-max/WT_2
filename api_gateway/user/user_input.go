package user

import (
	"time"
)

type Reps struct {
	Weight int `json:"weight"`
	Reps   int `json:"reps"`
}

type Workout struct {
	ExerciseId int    `json:"exercise_id"`
	RepsWeight []Reps `json:"tracker"`
}

type Tracker struct {
	Workout []Workout `json:"workout"`
}

type Signup struct {
	Name     string  `json:"name"`
	Email    string  `json:"email"`
	Password string  `json:"password"`
	Role     *string `json:"role"`
}

type Login struct {
	Message  string `json:"message"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UUIDReader struct {
	UUID string `json:"uuid"`
}

type ExerciseName struct {
	Name string `json:"exercise_name"`
}

type Exercise struct {
	Id        string       `json:"id,omitempty"`
	Name      string    `json:"name"`
	BodyPart  string    `json:"body_part"`
	Equipment string    `json:"equipment"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateExerResp struct{
	Messsage string `json:"message"`
	Exercise Exercise `json:"exercise"`
}

type Plan struct {
	PlanName  string   `json:"plan_name"`
	Exercises []string `json:"exercises"`
}

type PlanName struct {
	PlanName string `json:"plan_name"`
}

type ChangePass struct {
	OldPass string `json:"old_password"`
	NewPass string `json:"new_password"`
}
type ChangeEmail struct {
	OldEmail string `json:"old_email"`
	NewEmail string `json:"new_email"`
}
