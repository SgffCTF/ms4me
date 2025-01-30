package gamedto

import "errors"

var (
	ErrEmptyRequest = errors.New("request is empty")
)

type UpdateGameRequest struct {
	Title string `json:"title"`
	Mines int    `json:"mines"`
	Rows  int    `json:"rows"`
	Cols  int    `json:"cols"`
}

func (r *UpdateGameRequest) Validate() error {
	if r.Title == "" && r.Mines <= 0 && r.Rows <= 0 && r.Cols <= 0 {
		return ErrEmptyRequest
	}
	return nil
}
