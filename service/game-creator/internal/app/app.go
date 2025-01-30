package app

import (
	"context"
	"fmt"
	"game-creator/internal/config"
	grpcclient "game-creator/internal/grpc/client"
	gamehandlers "game-creator/internal/http/handlers"
	"game-creator/internal/http/middlewares"
	"game-creator/internal/services/game"
	"game-creator/internal/storage/postgres"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type App struct {
	log        *slog.Logger
	db         *postgres.Storage
	ssoClient  *grpcclient.SSOClient
	httpServer *http.Server
}

func New(cfg *config.ApplicationConfig, db *postgres.Storage, ssoClient *grpcclient.SSOClient, log *slog.Logger) *App {
	app := &App{
		log:       log,
		db:        db,
		ssoClient: ssoClient,
	}

	gameService := game.New(db)
	gameRouter := gamehandlers.New(log, gameService)

	httpServer := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Handler:      app.SetupRouter(gameRouter),
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
		IdleTimeout:  cfg.IdleTimeout,
	}
	app.httpServer = httpServer

	return app
}

func (app *App) Run() {
	err := app.httpServer.ListenAndServe()
	if err != nil {
		app.ssoClient.Close()
		app.db.Stop()
		panic("HTTP server failed to start: " + err.Error())
	}
}

func (app *App) Stop() {
	err := app.httpServer.Shutdown(context.Background())
	app.db.Stop()
	if err != nil {
		panic("HTTP server failed to stop: " + err.Error())
	}
}

func (app *App) SetupRouter(gameRouter *gamehandlers.GameRouter) http.Handler {
	router := chi.NewRouter()

	mw := middlewares.New(app.log, app.ssoClient.AuthClient)

	router.Use(middleware.Recoverer)
	router.Use(middleware.RequestID)
	router.Use(middleware.URLFormat)
	router.Use(mw.Logger())

	_ = router.Group(func(r chi.Router) {
		r.Use(mw.Auth())
		r.Post("/api/v1/game", gameRouter.CreateGame())
		r.Get("/api/v1/game", gameRouter.GetGames())
		r.Get("/api/v1/game/{id}", gameRouter.GetGame())
		r.Put("/api/v1/game/{id}", gameRouter.UpdateGame())
		r.Delete("/api/v1/game/{id}", gameRouter.DeleteGame())
	})

	router.Get("/api/v1/health", gamehandlers.Health())

	return router
}
