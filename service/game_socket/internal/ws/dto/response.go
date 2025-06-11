package dto_ws

import (
	"encoding/json"
	"errors"
)

var (
	ErrPlayerNotInGame = Error(errors.New("Игрок не подключен к игре"), AuthEventType)
)

type EventType string

const (
	SendMessageEventType    EventType = "SEND_MESSAGE"
	ReceiveMessageEventType EventType = "RECEIVE_MESSAGE"
	CreateRoomEventType     EventType = "CREATE_ROOM"
	UpdateRoomEventType     EventType = "UPDATE_ROOM"
	DeleteRoomEventType     EventType = "DELETE_ROOM"
	JoinRoomEventType       EventType = "JOIN_ROOM"
	ExitRoomEventType       EventType = "EXIT_ROOM"
	AuthEventType           EventType = "AUTH"

	StartGameEventType EventType = "START_GAME"
	ClickGameEventType EventType = "OPEN_CELL"
)

type Response struct {
	Status    string          `json:"status"`
	EventType EventType       `json:"event_type,omitempty"`
	Payload   json.RawMessage `json:"payload,omitempty"`
	Message   string          `json:"message,omitempty"`
	Error     string          `json:"error,omitempty"`
}

const (
	StatusOK    = "OK"
	StatusError = "Error"
)

func OK(msg string, et EventType) *Response {
	return &Response{Status: StatusOK, Message: msg, EventType: et}
}

func Error(err error, et EventType) *Response {
	return &Response{Status: StatusError, Error: err.Error(), EventType: et}
}

func (r *Response) Serialize() []byte {
	data, _ := json.Marshal(r)
	return data
}
