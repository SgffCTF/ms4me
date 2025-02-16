package ws

import (
	"context"
	"encoding/json"
	"errors"
	ssov1 "ms4me/game/pkg/grpc/sso"

	"google.golang.org/grpc/status"
)

type UserInfo interface {
	GetGameIDByUserID(ctx context.Context, userID int64) (string, error)
}

type Credentials struct {
	Token string `json:"token"`
}

type User struct {
	ID       int64
	Username string
}

type contextKey string

const userContextKey contextKey = "user"
const roomContextKey contextKey = "room"

// auth middleware validate token
func (s *Server) auth(client *Client) error {
	const op = "ws.auth"

	msg, err := s.read(client.conn)
	if err != nil {
		return ErrRead
	}

	var creds Credentials
	if err := json.Unmarshal(msg, &creds); err != nil {
		return ErrUnmarshalJSON
	}

	response, err := s.authClient.VerifyToken(client.ctx, &ssov1.VerifyTokenRequest{Token: creds.Token})
	if err != nil {
		return errors.New(status.Convert(err).Message())
	}
	user := response.GetUser()
	if user == nil {
		return ErrAuthError
	}
	client.ctx = context.WithValue(client.ctx, userContextKey, &User{ID: user.Id, Username: user.Username})

	gameID, err := s.db.GetGameIDByUserID(client.ctx, user.Id)
	if err != nil {
		return err
	}
	client.ctx = context.WithValue(client.ctx, roomContextKey, gameID)

	return nil
}
