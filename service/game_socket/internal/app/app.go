package app

import (
	"context"
	"log/slog"
	"ms4me/game_socket/internal/config"
	ws "ms4me/game_socket/internal/ws/server"
	"net"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	"golang.org/x/net/websocket"
)

type App struct {
	log        *slog.Logger
	cfg        *config.AppConfig
	httpServer *http.Server
	wsSrv      *ws.Server
}

func New(log *slog.Logger, cfg *config.AppConfig, wsSrv *ws.Server) *App {
	app := &App{
		cfg:   cfg,
		log:   log,
		wsSrv: wsSrv,
	}
	app.httpServer = &http.Server{
		Addr:         net.JoinHostPort(app.cfg.Host, strconv.Itoa(app.cfg.Port)),
		Handler:      app.initRouter(),
		ReadTimeout:  app.cfg.Timeout,
		WriteTimeout: app.cfg.Timeout,
		IdleTimeout:  app.cfg.IdleTimeout,
	}
	return app
}

func (a *App) Run() {
	err := a.httpServer.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		panic("HTTP server failed to start: " + err.Error())
	}
}

func (a *App) Stop(ctx context.Context) {
	err := a.httpServer.Shutdown(ctx)
	if err != nil {
		panic("error stopping http server: " + err.Error())
	}
}

func (a *App) initRouter() *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.Recoverer)
	router.Use(middleware.RequestID)
	router.Use(middleware.URLFormat)
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   a.cfg.CORSOrigins,
		AllowedMethods:   a.cfg.CORSMethods,
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// health route
	router.Get("/api/v1/health", func(w http.ResponseWriter, r *http.Request) {
		render.Data(w, r, []byte("OK"))
	})

	router.Handle("/ws", websocket.Handler(a.wsSrv.Handle))
	router.Handle("/ws/{id}", websocket.Handler(a.wsSrv.Handle))

	return router
}
