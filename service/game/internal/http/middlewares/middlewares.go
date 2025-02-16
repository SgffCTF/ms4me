package middlewares

import (
	"context"
	"log/slog"
	ssov1 "ms4me/game/pkg/grpc/sso"

	"google.golang.org/grpc"
)

type AuthService interface {
	VerifyToken(ctx context.Context, in *ssov1.VerifyTokenRequest, opts ...grpc.CallOption) (*ssov1.VerifyTokenResponse, error)
}

type Middlewares struct {
	log     *slog.Logger
	authSrv AuthService
}

func New(log *slog.Logger, authSrv AuthService) *Middlewares {
	return &Middlewares{
		log:     log,
		authSrv: authSrv,
	}
}
