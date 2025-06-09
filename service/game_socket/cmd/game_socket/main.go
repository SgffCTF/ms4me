package main

import (
	"context"
	"log/slog"
	"ms4me/game_socket/internal/app"
	"ms4me/game_socket/internal/config"
	storage "ms4me/game_socket/internal/redis"
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

	redisCli, err := storage.New(appCtx, cfg.RedisConfig)
	if err != nil {
		panic("error connecting to redis: " + err.Error())
	}
	wsSrv := ws.New(log, cfg.AppConfig, redisCli)
	eventLoop := eventloop.New(log, wsSrv, redisCli)
	go eventLoop.EventLoop()
	application := app.New(log, cfg.AppConfig, wsSrv)

	log.Info("starting application", slog.Any("config", cfg))
	go application.Run()

	sign := make(chan os.Signal, 1)

	signal.Notify(sign, syscall.SIGTERM, syscall.SIGINT)
	stopSignal := <-sign

	log.Info("stopping application", slog.String("signal", stopSignal.String()))
	eventLoop.Stop()
	application.Stop(appCtx)
}
