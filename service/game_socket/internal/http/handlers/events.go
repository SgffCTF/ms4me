package handlers

import (
	"log/slog"
	"ms4me/game_socket/internal/http/dto"
	"ms4me/game_socket/internal/http/dto/response"
	"net/http"

	"github.com/go-chi/render"
)

func (h *Handlers) Events() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req dto.EventRequest
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			render.JSON(w, r, response.ErrInternalError)
			return
		}

		for _, ev := range req.Events {
			*h.eventQueue <- ev
			h.log.Info("event from game service got successfully", slog.Any("event", ev))
		}
	}
}
