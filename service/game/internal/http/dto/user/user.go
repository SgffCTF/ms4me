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
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type RegisterResponse struct {
	response.Response
	ID int64 `json:"id"`
}
