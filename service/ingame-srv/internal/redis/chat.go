package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"ms4me/game_socket/internal/models"
	"time"
)

var MESSAGE_LIFETIME time.Duration = 20 * time.Minute

func (rc *Redis) CreateMessage(ctx context.Context, gameID string, message []byte) error {
	key := fmt.Sprintf("chat:%s", gameID)

	err := rc.DB.RPush(ctx, key, message).Err()
	if err != nil {
		return err
	}

	rc.DB.Expire(ctx, key, MESSAGE_LIFETIME)

	return nil
}

func (rc *Redis) ReadMessages(ctx context.Context, gameID string) ([]*models.Message, error) {
	key := fmt.Sprintf("chat:%s", gameID)

	messagesBytes, err := rc.DB.LRange(ctx, key, 0, -1).Result()
	if err != nil {
		return nil, err
	}

	var messages []*models.Message
	for _, msgBytes := range messagesBytes {
		var msg models.Message
		err := json.Unmarshal([]byte(msgBytes), &msg)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal message: %w", err)
		}
		messages = append(messages, &msg)
	}

	return messages, nil
}

func (rc *Redis) DeleteChat(ctx context.Context, gameID string) error {
	key := fmt.Sprintf("chat:%s", gameID)
	return rc.DB.Del(ctx, key).Err()
}

func (rc *Redis) ChatExists(ctx context.Context, roomID string) (bool, error) {
	key := fmt.Sprintf("chat:%s", roomID)

	exists, err := rc.DB.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check key existence: %w", err)
	}
	if exists == 0 {
		return false, nil
	}

	return true, nil
}
