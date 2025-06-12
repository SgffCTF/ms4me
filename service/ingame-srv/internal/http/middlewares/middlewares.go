package middlewares

import (
	"log/slog"
	gameclient "ms4me/game_socket/pkg/game_client"
)

type Middlewares struct {
	log        *slog.Logger
	jwtSecret  []byte
	gameClient *gameclient.GameClient
}

func New(log *slog.Logger, jwtSecret []byte, gameClient *gameclient.GameClient) *Middlewares {
	return &Middlewares{
		log:        log,
		jwtSecret:  jwtSecret,
		gameClient: gameClient,
	}
}
