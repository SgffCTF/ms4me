package handlers

import (
	"log/slog"
	storage "ms4me/game_socket/internal/redis"
)

type Handlers struct {
	log   *slog.Logger
	redis *storage.Redis
}

func New(log *slog.Logger, redis *storage.Redis) *Handlers {
	return &Handlers{
		log:   log,
		redis: redis,
	}
}
