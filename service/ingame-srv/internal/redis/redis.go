package storage

import (
	"context"
	"errors"
	"ms4me/game_socket/internal/config"
	"strconv"
	"time"

	redisdb "github.com/redis/go-redis/v9"
)

var (
	ErrNil = errors.New("key not found")
)

type Redis struct {
	DB     *redisdb.Client
	msgTTL time.Duration
}

func New(ctx context.Context, cfg *config.RedisConfig, msgTTL time.Duration) (*Redis, error) {
	client := redisdb.NewClient(&redisdb.Options{
		Addr:     cfg.Host + ":" + strconv.Itoa(cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
		Username: cfg.Username,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return &Redis{client, msgTTL}, nil
}

func (r *Redis) Stop(ctx context.Context) error {
	cmd := r.DB.ShutdownSave(ctx)
	if err := cmd.Err(); err != nil {
		return err
	}
	return nil
}
