package storage

import (
	"context"
	"encoding/json"
	"ms4me/game_socket/internal/models"
)

func (r *Redis) PublishEvent(ctx context.Context, event models.Event) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}
	err = r.DB.Publish(ctx, PUBLIC_QUEUE, data).Err()
	if err != nil {
		return err
	}
	return nil
}
