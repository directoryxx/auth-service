package domain

import "time"

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"-"`
}

type (
	RegisterValidation struct {
		Name     string `json:"name" validate:"required"`
		Email    string `json:"email" validate:"required"`
		Username string `json:"username" validate:"required"`
		Password string `json:"password" validate:"required"`
	}

	LoginValidation struct {
		Username string `json:"username" validate:"required"`
		Password string `json:"password" validate:"required"`
	}
)

type UserResponseAuthService struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Profile User   `json:"profile"`
}

type PublishAuthLogout struct {
	Data   LogoutAction
	Action string
}

type LogoutAction struct {
	Uuid string
}

type PublishAuthLogin struct {
	Data   LoginAction
	Action string
}

type LoginAction struct {
	Uuid string
	User User
	Exp  time.Duration
}
