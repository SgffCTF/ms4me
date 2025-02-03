package grpcclient

import (
	"context"
	"errors"
	"fmt"
	"ms4me/game/internal/config"
	ssov1 "ms4me/game/pkg/grpc/sso"
	"net"
	"strconv"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	ErrConnect = errors.New("can't connect to sso")
)

type SSOClient struct {
	conn         *grpc.ClientConn
	AuthClient   ssov1.AuthClient
	HealthClient ssov1.HealthClient
}

func New(cfg *config.SSOConfig) *SSOClient {
	conn, err := grpc.NewClient(
		net.JoinHostPort(cfg.Host, strconv.Itoa(cfg.Port)),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		panic("failed to create grpc client: " + err.Error())
	}

	authClient, healthClient := ssov1.NewAuthClient(conn), ssov1.NewHealthClient(conn)
	ssoClient := &SSOClient{conn: conn, AuthClient: authClient, HealthClient: healthClient}
	if err := ssoClient.Ping(); err != nil {
		panic(err)
	}

	return ssoClient
}

func (c *SSOClient) Ping() error {
	const op = "grpc.client.Ping"
	response, err := c.HealthClient.Ping(context.Background(), &ssov1.Empty{})
	if err != nil {
		return ErrConnect
	}
	if response.Message != "OK" {
		return fmt.Errorf("%s: %w", op, ErrConnect)
	}
	return nil
}

func (c *SSOClient) Close() {
	err := c.conn.Close()
	if err != nil {
		panic(err)
	}
}
