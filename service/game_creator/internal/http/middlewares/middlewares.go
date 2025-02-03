package middlewares

import (
	"log/slog"
	ssov1 "ms4me/game_creator/pkg/grpc/sso"
)

type Middlewares struct {
	log        *slog.Logger
	authClient ssov1.AuthClient
}

func New(log *slog.Logger, authClient ssov1.AuthClient) *Middlewares {
	return &Middlewares{
		log:        log,
		authClient: authClient,
	}
}
