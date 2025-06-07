package models

import "encoding/json"

type EventType int

const (
	TypeStartGame EventType = iota
	TypeCreateGame
	TypeJoinGame
	TypeDeleteGame
)

type Event struct {
	Type    EventType       `json:"type"`
	UserID  int64           `json:"user_id"`
	GameID  string          `json:"game_id"`
	Payload json.RawMessage `json:"payload,omitempty"`
}
