package main

import (
	"context"
	"game-creator/internal/app"
	"game-creator/internal/config"
	"game-creator/internal/storage/postgres"
	grpcclient "game-creator/pkg/grpc/client"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/jacute/prettylogger"
)

func main() {
	cfg := config.MustParseConfig()
	log := slog.New(prettylogger.NewColoredHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	appContext := context.Background()

	db := postgres.New(appContext, cfg.DBConfig)

	ssoClient := grpcclient.New(cfg.SSOConfig)

	application := app.New(cfg.AppConfig, db, ssoClient, log)

	log.Info("Starting app", slog.Any("config", cfg))
	go application.Run()

	sign := make(chan os.Signal, 1)
	signal.Notify(sign, syscall.SIGTERM, syscall.SIGINT)

	stopSignal := <-sign
	log.Info("stopping app", slog.String("signal", stopSignal.String()))

	application.Stop()
}
