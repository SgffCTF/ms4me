package gamedto

import (
	"errors"
	"game-creator/internal/http/dto/response"
)

var (
	ErrEmptyTitle = errors.New("title is empty")
	ErrMines      = errors.New("mines should be > 0")
	ErrCols       = errors.New("cols should be > 0")
	ErrRows       = errors.New("rows should be > 0")
)

type CreateGameRequest struct {
	Title    string `json:"title"`
	Mines    int    `json:"mines"`
	Rows     int    `json:"rows"`
	Cols     int    `json:"cols"`
	IsPublic *bool  `json:"is_public,omitempty"`
}

type CreateGameResponse struct {
	response.Response
	ID string `json:"id"`
}

func (r *CreateGameRequest) Validate() error {
	if r.Title == "" {
		return ErrEmptyTitle
	}
	if r.Mines <= 0 {
		return ErrMines
	}
	if r.Rows <= 0 {
		return ErrRows
	}
	if r.Cols <= 0 {
		return ErrCols
	}
	if r.IsPublic == nil {
		value := true
		r.IsPublic = &value
	}
	return nil
}
