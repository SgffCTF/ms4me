package userdto

import (
	"ms4me/game/internal/http/dto/response"
	"ms4me/game/internal/http/middlewares"
)

type UserResponse struct {
	response.Response
	User *middlewares.User `json:"user"`
}

type RegisterRequest struct {
	Username string `json:"username" validate:"required,max=150,min=4"`
	Password string `json:"password" validate:"required,min=8"`
}

type RegisterResponse struct {
	response.Response
	ID int64 `json:"id"`
}
