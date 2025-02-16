package handlers

import (
	"log/slog"
	eventsdto "ms4me/game/internal/http/dto/events"
	"ms4me/game/internal/http/dto/response"
	"ms4me/game/internal/service/events"
	"net/http"

	"github.com/go-chi/render"
)

type Handler struct {
	log       *slog.Logger
	eventsSrv *events.EventsService
}

func New(log *slog.Logger, eventsSrv *events.EventsService) *Handler {
	return &Handler{
		log:       log,
		eventsSrv: eventsSrv,
	}
}

func (h *Handler) Events() http.HandlerFunc {
	const op = "handlers.Events"
	log := h.log.With(slog.String("op", op))
	return func(w http.ResponseWriter, r *http.Request) {
		var req eventsdto.EventsRequest
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			render.JSON(w, r, response.Error(err.Error()))
			return
		}
		for _, event := range req.Events {
			h.eventsSrv.ProcessEvent(&event)
		}
		log.Info("events processed successfully", slog.Int("events_len", len(req.Events)))
		render.JSON(w, r, response.OK())
	}
}
