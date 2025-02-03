package main

import (
	"log/slog"
	"ms4me/game/internal/app"
	"ms4me/game/internal/config"
	grpcclient "ms4me/game/pkg/grpc/client"
	"os"
	"os/signal"
	"syscall"

	"github.com/jacute/prettylogger"
)

func main() {
	cfg := config.MustParse()
	var log *slog.Logger
	if cfg.Debug {
		log = slog.New(prettylogger.NewColoredHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	} else {
		log = slog.New(prettylogger.NewColoredHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	}

	ssoClient := grpcclient.New(cfg.SSOConfig)
	if err := ssoClient.Ping(); err != nil {
		panic(err)
	}

	application := app.New(cfg, log, ssoClient.AuthClient)
	log.Info("starting application", slog.Any("config", cfg))
	go application.Run()

	sign := make(chan os.Signal, 1)
	signal.Notify(sign, syscall.SIGTERM, syscall.SIGINT)

	stopSignal := <-sign
	application.Stop()
	log.Info("application stopped", slog.String("signal", stopSignal.String()))
}
