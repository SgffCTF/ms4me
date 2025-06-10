package main

import (
	"context"
	"log/slog"
	"ms4me/game/internal/app"
	"ms4me/game/internal/config"
	handlers "ms4me/game/internal/http/handlers"
	"ms4me/game/internal/services/auth"
	"ms4me/game/internal/services/batcher"
	"ms4me/game/internal/services/game"
	"ms4me/game/internal/storage/postgres"
	"ms4me/game/internal/storage/redis"
	gameclient "ms4me/game/pkg/game_client"
	"os"
	"os/signal"
	"syscall"

	"github.com/jacute/prettylogger"
)

func main() {
	cfg := config.MustParseConfig()

	var log *slog.Logger
	if cfg.Env == "local" {
		log = slog.New(prettylogger.NewColoredHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	} else {
		log = slog.New(prettylogger.NewColoredHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	appContext := context.Background()

	db := postgres.New(appContext, cfg.DatabaseConfig)

	rdb := redis.New(appContext, cfg.RedisConfig)
	batcher := batcher.New(log, rdb)
	log.Info("Starting batcher")
	go batcher.Start()

	gameSocketClient := gameclient.New(cfg.GameSocketConfig)
	gameService := game.New(log, db, batcher, gameSocketClient)
	authSrv := auth.New(log, db, []byte(cfg.JwtSecret), cfg.JwtTTL)
	gameHandlers := handlers.New(log, gameService, authSrv, cfg)

	application := app.New(cfg.ApplicationConfig, db, log, gameHandlers)
	log.Info("Starting app", slog.Any("config", cfg))
	go application.Run()

	sign := make(chan os.Signal, 1)
	signal.Notify(sign, syscall.SIGTERM, syscall.SIGINT)

	stopSignal := <-sign
	log.Info("stopping app", slog.String("signal", stopSignal.String()))
	application.Stop()
	log.Info("stopping batcher", slog.String("signal", stopSignal.String()))
	batcher.Shutdown()
}
