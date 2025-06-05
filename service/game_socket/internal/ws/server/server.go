package ws

import (
	"context"
	"errors"
	"log/slog"
	"sync"

	"ms4me/game_socket/internal/config"
	"ms4me/game_socket/internal/models"

	"golang.org/x/net/websocket"
)

const BUF_SIZE = 4096

type Server struct {
	cfg       *config.AppConfig
	log       *slog.Logger
	clients   map[int64]*Client
	rooms     map[string]*Room
	roomsMu   sync.Mutex
	clientsMu sync.Mutex
}

type Client struct {
	ctx       context.Context
	conn      *websocket.Conn
	user      *models.User
	requestID string
}

type Room struct {
	players []*models.User
}

var (
	ErrRead          = errors.New("read error")
	ErrUnmarshalJSON = errors.New("unmarshal JSON error")

	ErrAuthError = errors.New("auth error")
)

func New(log *slog.Logger, cfg *config.AppConfig) *Server {
	s := &Server{
		log:       log,
		clients:   make(map[int64]*Client),
		clientsMu: sync.Mutex{},
		cfg:       cfg,
	}
	return s
}
