package dto

import "ms4me/game_socket/internal/models"

type EventRequest struct {
	Events []models.Event `json:"events"`
}
