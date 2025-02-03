package gamedto

import (
	"errors"
	"ms4me/game_creator/internal/http/dto/response"

	validator "github.com/go-playground/validator/v10"
)

var (
	ErrEmptyTitle = errors.New("title is empty")
	ErrMines      = errors.New("mines should be > 0")
	ErrCols       = errors.New("cols should be > 0")
	ErrRows       = errors.New("rows should be > 0")
)

type CreateGameRequest struct {
	Title    string `json:"title" validate:"required"`
	Rows     int    `json:"rows" validate:"required,min=8,max=20"`
	Cols     int    `json:"cols" validate:"required,min=8,max=20"`
	IsPublic *bool  `json:"is_public,omitempty"`
}

type CreateGameResponse struct {
	response.Response
	ID string `json:"id"`
}

func (r *CreateGameRequest) Validate() error {
	if r.IsPublic == nil {
		value := true
		r.IsPublic = &value
	}
	validate := validator.New()
	return validate.Struct(r)
}
