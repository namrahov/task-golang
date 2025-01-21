package model

type UserRegistrationDto struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthRequestDto struct {
	EmailOrNickname string `json:"emailOrNickname"`
	Password        string `json:"password"`
	RememberMe      bool   `json:"rememberMe"`
}
