package app

import (
	"context"
	"log/slog"
	"ms4me/game/internal/config"
	"ms4me/game/internal/http/handlers"
	"ms4me/game/internal/http/middlewares"
	cent "ms4me/game/internal/service/centrifuge"
	"ms4me/game/internal/service/events"
	"net"
	"net/http"
	"time"

	"github.com/centrifugal/centrifuge"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type App struct {
	Node       *centrifuge.Node
	cfg        *config.Config
	log        *slog.Logger
	httpServer *http.Server
	centSrv    *cent.CentrifugeService
	eventsSrv  *events.EventsService
	authClient cent.AuthService
}

func New(cfg *config.Config, log *slog.Logger, authClient cent.AuthService) *App {
	node, err := centrifuge.New(centrifuge.Config{})
	if err != nil {
		panic(err)
	}
	centSrv := cent.New(node, log, authClient, cfg.CentrifugoConfig)
	node.OnConnecting(centSrv.OnConnecting)
	node.OnConnect(centSrv.OnConnect)
	eventsService := events.New(log, centSrv)

	httpServer := &http.Server{
		Addr:         net.JoinHostPort(cfg.AppConfig.Host, cfg.AppConfig.Port),
		ReadTimeout:  time.Duration(cfg.AppConfig.Timeout) * time.Second,
		WriteTimeout: time.Duration(cfg.AppConfig.Timeout) * time.Second,
		IdleTimeout:  time.Duration(cfg.AppConfig.IdleTimeout) * time.Second,
	}

	return &App{
		Node:       node,
		cfg:        cfg,
		log:        log,
		httpServer: httpServer,
		centSrv:    centSrv,
		eventsSrv:  eventsService,
		authClient: authClient,
	}
}

func (a *App) Run() {
	if err := a.Node.Run(); err != nil {
		panic(err)
	}

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

func (a *App) setupRouter() *chi.Mux {
	router := chi.NewRouter()

	mw := middlewares.New(a.log)

	router.Use(middleware.Recoverer)
	router.Use(middleware.RequestID)
	router.Use(middleware.URLFormat)
	router.Use(mw.Logger())

	handlers := handlers.New(a.log, a.eventsSrv)
	wsHandler := centrifuge.NewWebsocketHandler(a.Node, centrifuge.WebsocketConfig{
		ReadBufferSize:     1024,
		UseWriteBufferPool: true,
	})
	router.Handle("/connection/websocket", wsHandler)

	_ = router.Group(func(r chi.Router) {
		r.Post("/api/v1/events", handlers.Events())
	})

	return router
}
