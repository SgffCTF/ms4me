package main

import (
	"context"
	"log/slog"
	"ms4me/game_socket/internal/app"
	"ms4me/game_socket/internal/config"
	"ms4me/game_socket/internal/http/handlers"
	"ms4me/game_socket/internal/models"
	"ms4me/game_socket/internal/ws/eventloop"
	ws "ms4me/game_socket/internal/ws/server"
	"os"
	"os/signal"
	"syscall"

	"github.com/jacute/prettylogger"
)

const QUEUE_LEN = 50

func main() {
	appCtx := context.Background()
	cfg := config.MustParseConfig()

	var level slog.Level
	if cfg.Env == "local" {
		level = slog.LevelDebug
	} else {
		level = slog.LevelInfo
	}
	log := slog.New(prettylogger.NewColoredHandler(os.Stdout, &slog.HandlerOptions{Level: level}))

	eventsQueue := make(chan models.Event, QUEUE_LEN)
	wsSrv := ws.New(log, cfg.AppConfig)
	eventLoop := eventloop.New(log, &eventsQueue, wsSrv)
	go eventLoop.EventLoop()
	h := handlers.New(log, &eventsQueue)
	application := app.New(log, cfg.AppConfig, wsSrv, h)

	log.Info("starting application", slog.Any("config", cfg))
	go application.Run()

	sign := make(chan os.Signal, 1)

	signal.Notify(sign, syscall.SIGTERM, syscall.SIGINT)
	stopSignal := <-sign

	log.Info("stopping application", slog.String("signal", stopSignal.String()))
	eventLoop.Stop()
	application.Stop(appCtx)
}
