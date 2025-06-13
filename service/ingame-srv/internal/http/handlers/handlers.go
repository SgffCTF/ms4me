package handlers

import (
	"log/slog"
	storage "ms4me/game_socket/internal/redis"
	ws "ms4me/game_socket/internal/ws/server"
	gameclient "ms4me/game_socket/pkg/game_client"
)

type Handlers struct {
	log        *slog.Logger
	redis      *storage.Redis
	wsSrv      *ws.Server
	gameClient *gameclient.GameClient
}

func New(
	log *slog.Logger,
	redis *storage.Redis,
	wsSrv *ws.Server,
	gc *gameclient.GameClient,
) *Handlers {
	return &Handlers{
		log:        log,
		redis:      redis,
		wsSrv:      wsSrv,
		gameClient: gc,
	}
}
