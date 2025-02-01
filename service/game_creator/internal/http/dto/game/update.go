package gamedto

import (
	"errors"

	validator "github.com/go-playground/validator/v10"
)

var (
	ErrEmptyRequest = errors.New("request is empty")
)

type UpdateGameRequest struct {
	Title    string `json:"title"`
	Rows     int    `json:"rows" validate:"omitempty,min=8,max=20"`
	Cols     int    `json:"cols" validate:"omitempty,min=8,max=20"`
	IsPublic *bool  `json:"is_public,omitempty"`
}

func (r *UpdateGameRequest) Validate() error {
	if r.IsPublic == nil {
		value := true
		r.IsPublic = &value
	}
	if r.Title == "" && r.Rows <= 0 && r.Cols <= 0 {
		return ErrEmptyRequest
	}
	validate := validator.New()
	return validate.Struct(r)
}
