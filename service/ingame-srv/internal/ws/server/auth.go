package ws

import (
	"context"
	"encoding/json"
	"log/slog"
	"ms4me/game_socket/internal/models"
	"ms4me/game_socket/pkg/lib/jwt"

	"github.com/jacute/prettylogger"
	"golang.org/x/net/websocket"
)

type Credentials struct {
	Token string `json:"token"`
}

type contextKey string

const userContextKey contextKey = "user"
const roomContextKey contextKey = "room"

// auth middleware validate token
func (s *Server) auth(ctx context.Context, conn *websocket.Conn) (*models.User, error) {
	const op = "ws.auth"
	log := s.log.With(slog.String("op", op))

	msg, err := s.read(conn)
	if err != nil {
		return nil, ErrRead
	}

	var creds Credentials
	if err := json.Unmarshal(msg, &creds); err != nil {
		return nil, ErrUnmarshalJSON
	}

	claims, err := jwt.VerifyToken(creds.Token, []byte(s.cfg.JwtSecret))
	if err != nil {
		log.Info("error verifying token", prettylogger.Err(err))
		return nil, ErrAuthError
	}

	userID, ok := claims["user_id"].(float64)
	if !ok {
		return nil, ErrAuthError
	}
	username, ok := claims["username"].(string)

	return &models.User{
		ID:       int64(userID),
		Username: username,
	}, nil
}
