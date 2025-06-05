package gameclient

import (
	"ms4me/game/internal/models"
)

type EventsRequest struct {
	Events []models.Event `json:"events"`
}
