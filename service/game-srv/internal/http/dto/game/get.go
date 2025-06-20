package gamedto

import (
	"errors"
	"ms4me/game/internal/http/dto/response"
	"ms4me/game/internal/models"
	"net/url"
	"strconv"
)

var (
	ErrPage  = errors.New("page should be number > 0")
	ErrLimit = errors.New("limit should be number > 0")
)

type GetGamesRequest struct {
	Query  string
	Status string
	Page   int
	Limit  int
}

type GetGamesResponse struct {
	response.Response
	Games []*models.Game `json:"games"`
}

type GetGameResponse struct {
	response.Response
	Game *models.GameDetails `json:"game"`
}

func (ggr *GetGamesRequest) Render(values url.Values) error {
	if values.Has("query") {
		ggr.Query = values.Get("query")
	}
	if values.Has("status") {
		ggr.Status = values.Get("status")
	}
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

type GameStatusResponse struct {
	response.Response
	Result string `json:"result"`
}

type CloseGameRequest struct {
	WinnerID int64 `json:"winner_id"`
}

type GetCongratulationResponse struct {
	response.Response
	Congratulation string `json:"congratulation"`
}
