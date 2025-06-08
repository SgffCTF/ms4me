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
	Payload  json.RawMessage `json:"payload,omitempty"`
}
