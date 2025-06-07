package dto_ws

import (
	"encoding/json"
)

type EventType string

const (
	SendMessageEventType    EventType = "SEND_MESSAGE"
	ReceiveMessageEventType EventType = "RECEIVE_MESSAGE"
	CreateRoomEventType     EventType = "CREATE_ROOM"
	DeleteRoomEventType     EventType = "DELETE_ROOM"
	AuthEventType           EventType = "AUTH"
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
