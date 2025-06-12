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
