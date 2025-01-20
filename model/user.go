package model

type UserRegistrationDto struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
