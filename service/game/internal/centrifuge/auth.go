package service

import (
	"context"
	"errors"
	ssov1 "game/pkg/grpc/sso"

	"google.golang.org/grpc/status"
)

type User struct {
	ID       int64
	Username string
}

// auth middleware validate token
func (cs *CentrifugeService) auth(ctx context.Context, token string) (*User, error) {
	response, err := cs.authClient.VerifyToken(ctx, &ssov1.VerifyTokenRequest{Token: token})
	if err != nil {
		return nil, errors.New(status.Convert(err).Message())
	}
	return &User{ID: response.User.Id, Username: response.User.Username}, nil
}
