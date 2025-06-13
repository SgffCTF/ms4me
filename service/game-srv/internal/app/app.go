package app

import (
	"context"
	"fmt"
	"log/slog"
	"ms4me/game/internal/config"
	"ms4me/game/internal/http/handlers"
	"ms4me/game/internal/http/middlewares"
	"ms4me/game/internal/storage/postgres"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

type App struct {
	log        *slog.Logger
	db         *postgres.Storage
	cfg        *config.ApplicationConfig
	httpServer *http.Server
}

func New(cfg *config.ApplicationConfig, db *postgres.Storage, log *slog.Logger, h *handlers.GameHandlers) *App {
	app := &App{
		log: log,
		db:  db,
		cfg: cfg,
	}

	httpServer := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Handler:      app.SetupRouter(h),
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
		a.db.Stop()
		panic("HTTP server failed to start: " + err.Error())
	}
}

func (a *App) Stop() {
	a.db.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := a.httpServer.Shutdown(ctx)
	if err != nil {
		panic("HTTP server shutdown error: " + err.Error())
	}
}

func (a *App) SetupRouter(h *handlers.GameHandlers) http.Handler {
	router := chi.NewRouter()

	mw := middlewares.New(a.log, []byte(a.cfg.JwtSecret))

	router.Use(middleware.Recoverer)
	router.Use(middleware.RequestID)
	router.Use(middleware.URLFormat)
	router.Use(mw.Logger())
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   a.cfg.CORSOrigins,
		AllowedMethods:   a.cfg.CORSMethods,
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	router.Route("/api/v1/user", func(r chi.Router) {
		r.Post("/", h.Register())
		r.Get("/", mw.Auth()(h.User()).ServeHTTP)
		r.Post("/login", h.Login())
		r.Get("/logout", h.Logout())
		r.Get("/game", mw.Auth()(h.GetMyGames()).ServeHTTP)
	})

	router.Route("/api/v1/game", func(gameRouter chi.Router) {
		gameRouter.Use(mw.Auth())
		gameRouter.Post("/", h.CreateGame())
		gameRouter.Get("/", h.GetGames())
		gameRouter.Get("/{id}", h.GetGame())
		gameRouter.Put("/{id}", h.UpdateGame())
		gameRouter.Delete("/{id}", h.DeleteGame())
		gameRouter.Post("/{id}/start", h.StartGame())
		gameRouter.Post("/{id}/enter", h.EnterGame())
		gameRouter.Post("/{id}/exit", h.ExitGame())

		gameRouter.Post("/{id}/field/cell", h.OpenCell())
	})

	router.Route("/api/v1/internal", func(r chi.Router) {
		r.Get("/game/{id}/started", h.GameStarted())
		r.Post("/game/{id}/close", h.CloseGame())
	})

	router.Get("/api/v1/health", handlers.Health())

	return router
}
