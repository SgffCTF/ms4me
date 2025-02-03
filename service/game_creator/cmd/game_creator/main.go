package main

import (
	"context"
	"log/slog"
	"ms4me/game_creator/internal/app"
	"ms4me/game_creator/internal/config"
	gamehandlers "ms4me/game_creator/internal/http/handlers"
	"ms4me/game_creator/internal/services/batcher"
	"ms4me/game_creator/internal/services/game"
	"ms4me/game_creator/internal/storage/postgres"
	grpcclient "ms4me/game_creator/pkg/grpc/client"
	gameclient "ms4me/game_creator/pkg/http/game"
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

	db := postgres.New(appContext, cfg.DBConfig)
	ssoClient := grpcclient.New(cfg.SSOConfig)
	gameClient := gameclient.New(cfg.GameConfig)

	batcher := batcher.New(log, gameClient)
	log.Info("Starting batcher")
	go batcher.Start()

	gameService := game.New(db, batcher)
	gameHandlers := gamehandlers.New(log, gameService)

	application := app.New(cfg.AppConfig, db, ssoClient, log, gameHandlers)
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
