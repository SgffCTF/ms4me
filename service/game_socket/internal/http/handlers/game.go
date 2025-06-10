package handlers

import (
	"encoding/json"
	"log/slog"
	"ms4me/game_socket/internal/game"
	"ms4me/game_socket/internal/http/dto"
	"ms4me/game_socket/internal/http/middlewares"
	"ms4me/game_socket/internal/models"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/jacute/prettylogger"
)

func (h *Handlers) Ready() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		id := chi.URLParamFromCtx(ctx, "id")
		participants, err := h.redis.GetClientsInChannel(ctx, id)
		if err != nil {
			h.log.Error("error getting ready users", prettylogger.Err(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		for _, participant := range participants {
			result := h.wsSrv.CheckConn(participant.ID)
			if !result {
				w.WriteHeader(http.StatusTooEarly) // игроки не готовы к началу игры
				return
			}
		}
		w.WriteHeader(http.StatusOK)
	}
}

func (h *Handlers) OpenCell() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.OpenCell"

		ctx := r.Context()
		user := ctx.Value(middlewares.UserContextKey).(*middlewares.User)
		log := slog.With(slog.String("op", op), slog.Int64("user_id", user.ID))

		var req dto.OpenCellRequest
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, dto.ErrBody)
			return
		}

		id := chi.URLParamFromCtx(ctx, "id")
		participant, err := h.redis.GetClientInChannel(ctx, id, user.ID)
		if err != nil {
			log.Error("error getting room participants", prettylogger.Err(err))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, dto.ErrInternalError)
			return
		}

		field := participant.Field
		if field == nil {
			field = game.CreateField(req.Row, req.Col)
		} else {
			field.OpenCell(req.Row, req.Col)
		}
		fieldMarshalled, err := json.Marshal(field)
		if err != nil {
			log.Error("error marshalling field", prettylogger.Err(err))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, dto.ErrInternalError)
			return
		}
		err = h.redis.PublishEvent(ctx, models.Event{
			Type:     models.TypeGameField,
			UserID:   user.ID,
			GameID:   id,
			IsPublic: false,
			Payload:  fieldMarshalled,
		})
		if err != nil {
			h.log.Error("error publishing event", prettylogger.Err(err))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, dto.ErrInternalError)
			return
		}
	}
}
