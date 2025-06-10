package dto

const (
	StatusOK    = "OK"
	StatusError = "Error"
)

var (
	ErrInternalError = Error("Внутренняя ошибка")
	ErrBody          = Error("Ошибка формата тела запроса")
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
