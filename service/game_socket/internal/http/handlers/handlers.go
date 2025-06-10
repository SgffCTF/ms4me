package handlers

import (
	"log/slog"
	storage "ms4me/game_socket/internal/redis"
	ws "ms4me/game_socket/internal/ws/server"
)

type Handlers struct {
	log   *slog.Logger
	redis *storage.Redis
	wsSrv *ws.Server
}

func New(log *slog.Logger, redis *storage.Redis, wsSrv *ws.Server) *Handlers {
	return &Handlers{
		log:   log,
		redis: redis,
		wsSrv: wsSrv,
	}
}
