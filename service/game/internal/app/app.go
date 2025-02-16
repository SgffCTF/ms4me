package app

import (
	"context"
	"log/slog"
	"ms4me/game/internal/config"
	"ms4me/game/internal/http/handlers"
	ws "ms4me/game/internal/websocket"
	"net"
	"net/http"
	"time"

	"golang.org/x/net/websocket"
)

type App struct {
	wsServer   *ws.Server
	httpServer *http.Server
	handler    *handlers.Handler
	cfg        *config.Config
	log        *slog.Logger
}

func New(cfg *config.Config, log *slog.Logger, wsServer *ws.Server, handler *handlers.Handler) *App {
	httpServer := &http.Server{
		Addr:         net.JoinHostPort(cfg.AppConfig.Host, cfg.AppConfig.Port),
		ReadTimeout:  time.Duration(cfg.AppConfig.Timeout) * time.Second,
		WriteTimeout: time.Duration(cfg.AppConfig.Timeout) * time.Second,
		IdleTimeout:  time.Duration(cfg.AppConfig.IdleTimeout) * time.Second,
	}

	return &App{
		wsServer:   wsServer,
		cfg:        cfg,
		log:        log,
		httpServer: httpServer,
		handler:    handler,
	}
}

func (a *App) Run() {
	a.httpServer.Handler = a.setupRouter()

	if err := a.httpServer.ListenAndServe(); err != nil {
		panic(err)
	}
}

func (a *App) Stop() {
	err := a.httpServer.Shutdown(context.Background())
	if err != nil {
		panic("HTTP server failed to stop: " + err.Error())
	}
}

func (a *App) setupRouter() *http.ServeMux {
	router := http.NewServeMux()

	// mw := middlewares.New(a.log, a.authClient)

	// router.Use(middleware.Recoverer)
	// router.Use(middleware.RequestID)
	// router.Use(middleware.URLFormat)
	// router.Use(mw.Logger())

	// handlers := handlers.New(a.log, a.eventsSrv)

	// authMiddleware := chi.Chain(mw.Auth())
	router.Handle("/ws", websocket.Handler(a.wsServer.Handle))

	// _ = router.Group(func(r chi.Router) {
	// 	r.Post("/api/v1/events", handlers.Events())
	// })

	return router
}
