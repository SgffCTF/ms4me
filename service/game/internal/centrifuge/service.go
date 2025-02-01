package service

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/centrifugal/centrifuge"
	"github.com/jacute/prettylogger"
	"google.golang.org/grpc"

	"game/internal/config"
	"game/internal/storage/redis"
	ssov1 "game/pkg/grpc/sso"
)

type AuthService interface {
	VerifyToken(ctx context.Context, in *ssov1.VerifyTokenRequest, opts ...grpc.CallOption) (*ssov1.VerifyTokenResponse, error)
}

type Credentials struct {
	Token string `json:"token"`
}

type CentrifugeService struct {
	node       *centrifuge.Node
	log        *slog.Logger
	authClient AuthService
	cfg        *config.CentrifugoConfig
}

func New(node *centrifuge.Node, log *slog.Logger, authClient AuthService, cfg *config.CentrifugoConfig) *CentrifugeService {
	return &CentrifugeService{node: node, log: log, authClient: authClient, cfg: cfg}
}

func (cs *CentrifugeService) OnConnecting(ctx context.Context, c centrifuge.ConnectEvent) (centrifuge.ConnectReply, error) {
	const op = "centrifuge.OnConnecting"
	log := cs.log.With(slog.String("op", op))

	log.Debug(
		"client onconnecting",
		slog.String("transport", c.Transport.Name()),
		slog.String("transport_proto", string(c.Transport.Protocol())),
	)

	user, err := cs.auth(ctx, c.Token)
	if err != nil {
		log.Warn("error authorizing", prettylogger.Err(err))
		return centrifuge.ConnectReply{}, &centrifuge.Error{
			Code:    101,
			Message: err.Error(),
		}
	}
	userID := strconv.FormatInt(user.ID, 10)

	channel := redis.GetChannel(userID)
	exp := time.Now().UTC().Add(cs.cfg.ExpTime).Unix()

	cs.log.Debug("client connecting", slog.Int64("user_id", user.ID))

	fmt.Println(userID, channel)
	if err := cs.node.Subscribe(userID, channel); err != nil {
		return centrifuge.ConnectReply{}, nil
	}

	return centrifuge.ConnectReply{
		Credentials: &centrifuge.Credentials{
			UserID:   userID,
			ExpireAt: exp,
		},
		Data: []byte("{}"),
		PingPongConfig: &centrifuge.PingPongConfig{
			PingInterval: cs.cfg.PingInterval,
			PongTimeout:  cs.cfg.PongTimeout,
		},
		Subscriptions: map[string]centrifuge.SubscribeOptions{},
	}, nil
}

func (cs *CentrifugeService) OnConnect(client *centrifuge.Client) {
	const op = "centrifuge.OnConnect"
	log := cs.log.With(slog.String("op", op))

	transportName := client.Transport().Name()
	transportProto := client.Transport().Protocol()
	log.Info("client connected", slog.String("transport_name", transportName), slog.String("transport_proto", string(transportProto)))

	client.OnSubscribe(func(e centrifuge.SubscribeEvent, cb centrifuge.SubscribeCallback) {
		log.Info("client subscribes on channel", slog.String("channel", e.Channel))
		cb(centrifuge.SubscribeReply{}, nil)
	})

	client.OnPublish(func(e centrifuge.PublishEvent, cb centrifuge.PublishCallback) {
		log.Debug("client publishes message into channel", slog.String("channel", e.Channel), slog.String("message", string(e.Data)))
		cb(centrifuge.PublishReply{}, nil)
	})

	client.OnDisconnect(func(e centrifuge.DisconnectEvent) {
		log.Info("client disconnected", slog.String("reason", e.Reason))
	})
}
