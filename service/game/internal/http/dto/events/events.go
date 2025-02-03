package eventsdto

import "encoding/json"

type EventType int

const (
	StartGame EventType = iota
	CreateGame
	JoinGame
)

type Event struct {
	Type    EventType       `json:"type"`
	UserID  int64           `json:"user_id"`
	GameID  string          `json:"game_id"`
	Payload json.RawMessage `json:"payload,omitempty"`
}

type EventsRequest struct {
	Events []Event `json:"events"`
}
