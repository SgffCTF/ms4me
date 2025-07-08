package main

import (
	"context"
	"log/slog"
	"ms4me/cleaner/internal/app"
	"ms4me/cleaner/internal/config"
	"ms4me/cleaner/internal/storage/postgres"
	"os"
	"os/signal"
	"syscall"

	"github.com/jacute/prettylogger"
)

func main() {
	log := slog.New(prettylogger.NewJsonHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	cfg := config.MustParseConfig()
	storage := postgres.New(context.Background(), cfg.DatabaseConfig)

	application := app.New(log, cfg, storage)
	log.Info("running app")
	go application.Start()
	sign := make(chan os.Signal, 1)
	signal.Notify(sign, syscall.SIGTERM, syscall.SIGINT)

	stopSignal := <-sign
	log.Info("stopping app", slog.String("signal", stopSignal.String()))
	application.Stop()
}
