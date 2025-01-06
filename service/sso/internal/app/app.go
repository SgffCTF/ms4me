package app

import (
	"log/slog"
	"ms4me/sso/internal/app/grpcapp"
	"ms4me/sso/internal/config"
	"ms4me/sso/internal/database/postgres"
	"ms4me/sso/internal/services/auth"
)

type App struct {
	GRPCApp *grpcapp.App
	DB      *postgres.Storage
}

func New(cfg *config.Config, log *slog.Logger) *App {
	db := postgres.New(cfg.DBConfig)

	authService := auth.New(log, db, db, []byte(cfg.AppConfig.JwtSecret), cfg.AppConfig.JwtTTL)

	return &App{
		GRPCApp: grpcapp.New(log, authService, cfg.AppConfig.Host, cfg.AppConfig.Port),
		DB:      db,
	}
}
