package middlewares

import (
	ssov1 "game-creator/internal/grpc/sso"
	"log/slog"
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
