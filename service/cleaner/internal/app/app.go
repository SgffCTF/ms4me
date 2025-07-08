package app

import (
	"context"
	"log/slog"
	"ms4me/cleaner/internal/config"
	"ms4me/cleaner/internal/storage/postgres"
	"time"

	"github.com/jacute/prettylogger"
)

type Storage interface {
	DeleteGamesBefore(ctx context.Context, t time.Time) (int64, error)
}

type App struct {
	log    *slog.Logger
	cfg    *config.Config
	db     Storage
	stopCh chan struct{}
}

func New(log *slog.Logger, cfg *config.Config, db *postgres.Storage) *App {
	return &App{
		log:    log,
		cfg:    cfg,
		db:     db,
		stopCh: make(chan struct{}),
	}
}

func (a *App) Start() {
	const op = "app.Start"
	log := a.log.With(slog.String("op", op))

	appCtx, cancel := context.WithCancel(context.Background())
	t := time.NewTicker(a.cfg.CleanTimeout)
	for {
		select {
		case <-t.C:
			rowsAffected, err := a.db.DeleteGamesBefore(appCtx, time.Now().UTC().Add(-a.cfg.CleanBefore))
			if err != nil {
				log.Error("error deleting games before", prettylogger.Err(err))
				continue
			}
			log.Info("old games deleted", slog.Int64("count", rowsAffected))
		case <-a.stopCh:
			cancel()
			t.Stop()
			break
		}
	}
}

func (a *App) Stop() {
	a.stopCh <- struct{}{}
}
