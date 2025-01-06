package suite

import (
	"context"
	"log/slog"
	"ms4me/sso/internal/app"
	"ms4me/sso/internal/config"
	ssov1 "ms4me/sso/internal/grpc/proto/sso"
	"net"
	"strconv"
	"testing"
	"time"

	"github.com/jacute/prettylogger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const TIMEOUT = 25 * time.Second
const configPath = "./test_config.yml"

type Suite struct {
	Cfg        *config.Config
	App        *app.App
	AuthClient ssov1.AuthClient
}

func New(t *testing.T) (context.Context, *Suite) {
	cfg, err := config.ParseConfig(configPath)
	if err != nil {
		panic(err)
	}
	log := slog.New(prettylogger.NewDiscardHandler())

	ctx, cancelCtx := context.WithTimeout(context.Background(), TIMEOUT)

	var cc *grpc.ClientConn

	t.Cleanup(func() {
		t.Helper()
		cancelCtx()
		cc.Close()
	})

	app := app.New(cfg, log)
	go app.GRPCApp.MustRun()

	cc, err = grpc.NewClient(
		net.JoinHostPort(cfg.AppConfig.Host, strconv.Itoa(cfg.AppConfig.Port)),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		panic("error connecting to server: " + err.Error())
	}

	return ctx, &Suite{
		Cfg:        cfg,
		App:        app,
		AuthClient: ssov1.NewAuthClient(cc),
	}
}
