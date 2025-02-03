package app

import (
	"context"
	"fmt"
	"log/slog"
	"ms4me/game_creator/internal/config"
	gamehandlers "ms4me/game_creator/internal/http/handlers"
	"ms4me/game_creator/internal/http/middlewares"
	"ms4me/game_creator/internal/storage/postgres"
	grpcclient "ms4me/game_creator/pkg/grpc/client"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type App struct {
	log        *slog.Logger
	db         *postgres.Storage
	ssoClient  *grpcclient.SSOClient
	httpServer *http.Server
}

func New(cfg *config.ApplicationConfig, db *postgres.Storage, ssoClient *grpcclient.SSOClient, log *slog.Logger, gameRouter *gamehandlers.GameHandlers) *App {
	app := &App{
		log:       log,
		db:        db,
		ssoClient: ssoClient,
	}

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

func (a *App) Run() {
	err := a.httpServer.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		a.ssoClient.Close()
		a.db.Stop()
		panic("HTTP server failed to start: " + err.Error())
	}
}

func (a *App) Stop() {
	a.ssoClient.Close()
	a.db.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := a.httpServer.Shutdown(ctx)
	if err != nil {
		panic("HTTP server shutdown error: " + err.Error())
	}
}

func (a *App) SetupRouter(gameRouter *gamehandlers.GameHandlers) http.Handler {
	router := chi.NewRouter()

	mw := middlewares.New(a.log, a.ssoClient.AuthClient)

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
		r.Post("/api/v1/game/{id}/start", gameRouter.StartGame())
		r.Post("/api/v1/game/{id}/enter", gameRouter.EnterGame())
		r.Post("/api/v1/game/{id}/exit", gameRouter.ExitGame())
		// r.Route("/api/v1/game/{id}/reveal/{row}/{col}", func(r chi.Router) {
		// 	r.Get("/", gameRouter.RevealCell())
		// })
	})

	router.Get("/api/v1/health", gamehandlers.Health())

	return router
}
