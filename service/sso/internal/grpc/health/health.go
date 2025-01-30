package healthgrpc

import (
	"context"
	ssov1 "ms4me/sso/internal/grpc/proto/sso"

	"google.golang.org/grpc"
)

type serverAPI struct {
	ssov1.UnimplementedHealthServer
}

func Register(gRPC *grpc.Server) {
	ssov1.RegisterHealthServer(gRPC, &serverAPI{})
}

func (s *serverAPI) Ping(context.Context, *ssov1.Empty) (*ssov1.Pong, error) {
	return &ssov1.Pong{Message: "OK"}, nil
}
