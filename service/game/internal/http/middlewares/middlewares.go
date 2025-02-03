package middlewares

import (
	"log/slog"
)

type Middlewares struct {
	log *slog.Logger
}

func New(log *slog.Logger) *Middlewares {
	return &Middlewares{
		log: log,
	}
}
