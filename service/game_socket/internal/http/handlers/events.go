package handlers

import (
	"log/slog"
	"ms4me/game_socket/internal/http/dto"
	"ms4me/game_socket/internal/http/dto/response"
	"net/http"

	"github.com/go-chi/render"
	"github.com/jacute/prettylogger"
)

func (h *Handlers) Events() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req dto.EventRequest
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			render.JSON(w, r, response.ErrInternalError)
			return
		}

		for _, ev := range req.Events {
			err := h.redis.PublishEvent(r.Context(), ev)
			if err != nil {
				h.log.Error("error publishing event", prettylogger.Err(err))
				continue
			}
			h.log.Info("event from game service got successfully", slog.Any("event", ev))
		}
	}
}
