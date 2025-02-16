package ws

import (
	"encoding/json"
	"errors"
)

type Error struct {
	Err string `json:"error"`
}

var (
	ErrRead          = errors.New("read error")
	ErrUnmarshalJSON = errors.New("unmarshal JSON error")

	ErrAuthError = errors.New("auth error")
)

func NewError(err error) *Error {
	return &Error{
		Err: err.Error(),
	}
}

func (e *Error) Error() []byte {
	bytes, _ := json.Marshal(e)
	return bytes
}
