package handlers

import (
	"log/slog"
	"ms4me/game_socket/internal/models"
)

type Handlers struct {
	log        *slog.Logger
	eventQueue *chan models.Event
}

func New(log *slog.Logger, eventQueue *chan models.Event) *Handlers {
	return &Handlers{
		log:        log,
		eventQueue: eventQueue,
	}
}
