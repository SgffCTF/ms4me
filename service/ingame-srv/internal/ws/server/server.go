package ws

import (
	"context"
	"errors"
	"log/slog"
	"sync"

	"ms4me/game_socket/internal/config"
	"ms4me/game_socket/internal/models"
	storage "ms4me/game_socket/internal/redis"

	"golang.org/x/net/websocket"
)

const BUF_SIZE = 4096

type Server struct {
	cfg     *config.AppConfig
	log     *slog.Logger
	users   map[int64][]*Client
	usersMu sync.Mutex
	redis   *storage.Redis
}

type Client struct {
	ctx       context.Context
	conn      *websocket.Conn
	user      *models.User
	room      string
	requestID string
}

var (
	ErrRead          = errors.New("read error")
	ErrUnmarshalJSON = errors.New("unmarshal JSON error")

	ErrAuthError = errors.New("auth error")
)

func New(log *slog.Logger, cfg *config.AppConfig, redis *storage.Redis) *Server {
	s := &Server{
		log:     log,
		users:   make(map[int64][]*Client),
		usersMu: sync.Mutex{},
		cfg:     cfg,
		redis:   redis,
	}
	return s
}
