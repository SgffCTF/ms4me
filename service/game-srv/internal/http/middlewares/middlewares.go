package middlewares

import (
	"log/slog"
)

type Middlewares struct {
	log       *slog.Logger
	jwtSecret []byte
}

func New(log *slog.Logger, jwtSecret []byte) *Middlewares {
	return &Middlewares{
		log:       log,
		jwtSecret: jwtSecret,
	}
}
