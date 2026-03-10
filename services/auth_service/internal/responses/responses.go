package responses

import (
	"time"
)

type response struct {
	Message string `json:"message"`
}

var (
	SignUPSuccessFull = response{Message: "Successfully signed up"}
	LoginSuccessFull  = response{Message: "login Successfull"}
)

type SignUpResp struct {
	Name      string    `json:"name" validate:"min=2, max=20"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

type LoginResp struct {
	Message     string `json:"message"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	AccessToken string `json:"access_token"`
	UUID        string `json:"uuid"`
}
