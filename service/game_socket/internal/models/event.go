package models

import "encoding/json"

type EventType int

const (
	TypeStartGame EventType = iota
	TypeCreateGame
	TypeJoinGame
	TypeDeleteGame
	TypeUpdateGame
)

type Event struct {
	Type     EventType       `json:"type"`
	UserID   int64           `json:"user_id"`
	Username string          `json:"username,omitempty"`
	GameID   string          `json:"game_id"`
	IsPublic bool            `json:"is_public,omitempty"`
	Payload  json.RawMessage `json:"payload,omitempty"`
}

type CreateEvent struct {
	ID        string `json:"id"`
	OwnerID   int64  `json:"owner_id"`
	OwnerName string `json:"owner_name"`
}

type RoomParticipant struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	IsOwner  bool   `json:"is_owner"`
}
