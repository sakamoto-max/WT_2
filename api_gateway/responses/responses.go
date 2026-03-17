package responses

import "time"

type SignUpResp struct {
	Name      string    `json:"name" validate:"min=2, max=20"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}
