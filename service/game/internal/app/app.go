package app

import (
	"context"
	service "game/internal/centrifuge"
	"game/internal/config"
	ssov1 "game/pkg/grpc/sso"
	"log/slog"
	"net"
	"net/http"
	"time"

	"github.com/centrifugal/centrifuge"
	"google.golang.org/grpc"
)

type AuthService interface {
	VerifyToken(ctx context.Context, in *ssov1.VerifyTokenRequest, opts ...grpc.CallOption) (*ssov1.VerifyTokenResponse, error)
}

type App struct {
	Node       *centrifuge.Node
	Cfg        *config.Config
	Log        *slog.Logger
	httpServer *http.Server
	srv        *service.CentrifugeService
	authClient AuthService
}

func New(cfg *config.Config, log *slog.Logger, authClient AuthService) *App {

	node, err := centrifuge.New(centrifuge.Config{})
	if err != nil {
		panic(err)
	}
	srv := service.New(node, log, authClient, cfg.CentrifugoConfig)
	node.OnConnecting(srv.OnConnecting)
	node.OnConnect(srv.OnConnect)

	httpServer := &http.Server{
		Addr:         net.JoinHostPort(cfg.AppConfig.Host, cfg.AppConfig.Port),
		ReadTimeout:  time.Duration(cfg.AppConfig.Timeout) * time.Second,
		WriteTimeout: time.Duration(cfg.AppConfig.Timeout) * time.Second,
		IdleTimeout:  time.Duration(cfg.AppConfig.IdleTimeout) * time.Second,
	}

	return &App{
		Node:       node,
		Cfg:        cfg,
		Log:        log,
		httpServer: httpServer,
		srv:        srv,
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

func (a *App) setupRouter() *http.ServeMux {
	router := http.NewServeMux()
	wsHandler := centrifuge.NewWebsocketHandler(a.Node, centrifuge.WebsocketConfig{})
	router.Handle("/ws", wsHandler)

	return router
}
