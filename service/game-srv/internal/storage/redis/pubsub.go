package redis

import (
	"context"
	"encoding/json"
	"ms4me/game/internal/models"
)

const PUBLIC_QUEUE = "queue"

func (r *Redis) PublishEvents(ctx context.Context, events []models.Event) error {
	pipe := r.DB.Pipeline()
	for _, event := range events {
		data, err := json.Marshal(event)
		if err != nil {
			return err
		}
		pipe.Publish(ctx, PUBLIC_QUEUE, data)
	}
	_, err := pipe.Exec(ctx)
	return err
}

func (r *Redis) PublishEvent(ctx context.Context, event models.Event) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}
	return r.DB.Publish(ctx, PUBLIC_QUEUE, data).Err()
}
