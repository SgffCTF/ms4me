package main

import (
	"log/slog"
	"ms4me/sso/internal/app"
	"ms4me/sso/internal/config"
	"os"
	"os/signal"
	"syscall"

	"github.com/jacute/prettylogger"
)

func main() {
	cfg := config.MustParseConfig()
	log := slog.New(prettylogger.NewColoredHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	application := app.New(cfg, log)
	go application.GRPCApp.Run()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	sign := <-stop
	application.GRPCApp.Stop()

	log.Info("Application stopped", slog.String("signal", sign.String()))
}
