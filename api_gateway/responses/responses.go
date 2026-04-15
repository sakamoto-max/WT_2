package responses

import "time"

type SignUpResp struct {
	Name      string    `json:"name" validate:"min=2, max=20"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}


type TrackerConflict struct {
	RequestStatus string `json:"request_status"`
	Reason string `json:"reason"`
	Message string `json:"message"`
	ExerciseNames []string `json:"exercise_names"`
}