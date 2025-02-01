package response

import "errors"

const (
	StatusOK    = "OK"
	StatusError = "Error"
)

var (
	ErrInternalError = errors.New("internal error")
)

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

func OK() Response {
	return Response{Status: StatusOK}
}

func Error(msg string) Response {
	return Response{Status: StatusError, Error: msg}
}
