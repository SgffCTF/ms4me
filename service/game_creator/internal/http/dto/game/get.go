package gamedto

import (
	"errors"
	"game-creator/internal/http/dto/response"
	"game-creator/internal/models"
	"net/url"
	"strconv"
)

var (
	ErrPage  = errors.New("page should be number > 0")
	ErrLimit = errors.New("limit should be number > 0")
)

type GetGamesRequest struct {
	Page  int
	Limit int
}

type GetGamesResponse struct {
	response.Response
	Games []*models.Game `json:"games"`
}

func (ggr *GetGamesRequest) Render(values url.Values) error {
	if !values.Has("page") && !values.Has("limit") {
		return nil
	}

	page, err := strconv.Atoi(values.Get("page"))
	if err != nil || page <= 0 {
		return ErrPage
	}

	limit, err := strconv.Atoi(values.Get("limit"))
	if err != nil || limit <= 0 {
		return ErrLimit
	}
	ggr.Page, ggr.Limit = page, limit

	return nil
}
