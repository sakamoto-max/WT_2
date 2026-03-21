package models

import (
	"time"

	"github.com/google/uuid"
)

type Signup struct {
	Name     string `json:"name" validate:"min=2, max=20"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Login struct {
	Message  string `json:"message"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UUIDReader struct {
	UUID uuid.UUID `json:"uuid"`
}

type ClaimsContext string

var Claimskey ClaimsContext

type MqMsg struct {
	UserId int
	Action string
	Time   time.Time
	Email  string
}

type UserIdPayload struct{
	UserId int `json:"user_id"`
}