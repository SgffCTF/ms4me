package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"ms4me/game_socket/internal/models"
	"strconv"
)

const PUBLIC_QUEUE = "queue"

func (rc *Redis) AddClientToChannel(ctx context.Context, channel string, userID int64, meta *models.RoomParticipant) error {
	key := fmt.Sprintf("room:%s", channel)
	data, err := json.Marshal(meta)
	if err != nil {
		return err
	}
	return rc.DB.HSet(ctx, key, fmt.Sprintf("%d", userID), data).Err()
}

func (rc *Redis) RemoveClientFromChannel(ctx context.Context, channel string, userID int64) error {
	key := fmt.Sprintf("room:%s", channel)
	return rc.DB.HDel(ctx, key, fmt.Sprintf("%d", userID)).Err()
}

func (rc *Redis) DeleteChannel(ctx context.Context, channel string) error {
	key := fmt.Sprintf("room:%s", channel)
	return rc.DB.Del(ctx, key).Err()
}

func (rc *Redis) GetClientsInChannel(ctx context.Context, channel string) (map[string]*models.RoomParticipant, error) {
	key := fmt.Sprintf("room:%s", channel)
	result, err := rc.DB.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	clients := make(map[string]*models.RoomParticipant)
	for uid, raw := range result {
		var meta *models.RoomParticipant
		if err := json.Unmarshal([]byte(raw), &meta); err != nil {
			continue
		}
		clients[uid] = meta
	}
	return clients, nil
}

func (rc *Redis) GetUsersInChannel(ctx context.Context, channel string) ([]int, error) {
	participants, err := rc.GetClientsInChannel(ctx, channel)
	if err != nil {
		return nil, err
	}
	users := make([]int, 0)
	for userIDStr, _ := range participants {
		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			return nil, err
		}
		users = append(users, userID)
	}
	return users, nil
}

func (rc *Redis) GetClientInChannel(ctx context.Context, roomID string, userID int64) (*models.RoomParticipant, error) {
	key := fmt.Sprintf("room:%s", roomID)
	result, err := rc.DB.HGet(ctx, key, fmt.Sprintf("%d", userID)).Result()
	if err != nil {
		return nil, err
	}

	var roomParticipant models.RoomParticipant
	if err := json.Unmarshal([]byte(result), &roomParticipant); err != nil {
		return nil, err
	}
	return &roomParticipant, nil
}
