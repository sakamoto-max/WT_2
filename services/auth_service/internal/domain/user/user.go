package user

type SignUpPayload struct {
	Name string
	Email string
	Password string
	Role string
}

type LoginPayload struct {
	Email string
	Password string
}

type ChangePassPayload struct {
	UserId string
	OldPass string
	NewPass string

}

type ChangeEmailPayload struct {
	UserId string
	OldEmail string
	NewEmail string

}