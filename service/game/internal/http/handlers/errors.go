package handlers

import (
	"ms4me/game/internal/http/dto/response"
)

var (
	ErrInvalidBody       = response.Error("Неправильный запрос")
	ErrEmptyID           = response.Error("ID не должен быть пустым")
	ErrUserNotFound      = response.Error("Пользователь не найден")
	ErrIncorrectPassword = response.Error("Неверный пароль")
)
