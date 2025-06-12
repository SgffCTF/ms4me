package models

import "encoding/json"

type EventType int

const (
	TypeStartGame EventType = iota
	TypeCreateGame
	TypeJoinGame
	TypeDeleteGame
	TypeUpdateGame
	TypeExitGame

	TypeOpenCell
)

type Event struct {
	Type     EventType       `json:"type"`
	UserID   int64           `json:"user_id"`
	Username string          `json:"username,omitempty"`
	GameID   string          `json:"game_id"`
	IsPublic bool            `json:"is_public,omitempty"`
	Payload  json.RawMessage `json:"payload,omitempty"`
}
