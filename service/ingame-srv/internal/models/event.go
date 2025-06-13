package models

import (
	"encoding/json"
	"ms4me/game_socket/internal/service/game"
)

type EventType int

const (
	TypeStartGame EventType = iota
	TypeCreateGame
	TypeJoinGame
	TypeDeleteGame
	TypeUpdateGame
	TypeExitGame

	TypeClickGame
	TypeLoseGame
	TypeWinGame

	TypeNewMessage
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

type ClickEvent struct {
	ID       int64       `json:"id"`
	Username string      `json:"username"`
	IsOwner  bool        `json:"is_owner"`
	Field    *game.Field `json:"field"`
}

type LoseEvent struct {
	LoserID       int64  `json:"loser_id"`
	LoserUsername string `json:"loser_username"`
}

type WinEvent struct {
	WinnerID       int64  `json:"winner_id"`
	WinnerUsername string `json:"winner_username"`
}

type RoomParticipant struct {
	ID       int64       `json:"id"`
	Username string      `json:"username"`
	IsOwner  bool        `json:"is_owner"`
	Field    *game.Field `json:"field"`
}
