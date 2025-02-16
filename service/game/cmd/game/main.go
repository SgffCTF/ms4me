package main

import (
	"context"
	"log/slog"
	"ms4me/game/internal/app"
	"ms4me/game/internal/config"
	"ms4me/game/internal/http/handlers"
	"ms4me/game/internal/service/events"
	"ms4me/game/internal/storage/channel"
	"ms4me/game/internal/storage/postgres"
	ws "ms4me/game/internal/websocket"
	grpcclient "ms4me/game/pkg/grpc/client"
	"os"
	"os/signal"
	"syscall"

	"github.com/jacute/prettylogger"
)

func main() {
	cfg := config.MustParse()
	var log *slog.Logger
	if cfg.Env == "local" {
		log = slog.New(prettylogger.NewColoredHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	} else {
		log = slog.New(prettylogger.NewColoredHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	}

	// storages
	channelStorage := channel.New()
	db := postgres.New(context.Background(), cfg.DatabaseConfig)

	// services
	ssoClient := grpcclient.New(cfg.SSOConfig)
	if err := ssoClient.Ping(); err != nil {
		panic(err)
	}
	eventsService := events.New(log, channelStorage)

	// servers
	wsServer := ws.New(log, ssoClient.AuthClient, db)
	handler := handlers.New(log, eventsService)

	application := app.New(cfg, log, wsServer, handler)
	log.Info("starting application", slog.Any("config", cfg))
	go application.Run()

	sign := make(chan os.Signal, 1)
	signal.Notify(sign, syscall.SIGTERM, syscall.SIGINT)

	stopSignal := <-sign
	application.Stop()
	log.Info("application stopped", slog.String("signal", stopSignal.String()))
}
