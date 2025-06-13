package handlers

import (
	"encoding/json"
	"log/slog"
	"ms4me/game_socket/internal/http/dto"
	"ms4me/game_socket/internal/http/middlewares"
	"ms4me/game_socket/internal/models"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/jacute/prettylogger"
)

func (h *Handlers) CreateMessage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.CreateMessage"

		ctx := r.Context()
		user := ctx.Value(middlewares.UserContextKey).(*middlewares.User)
		w.Header().Set("Content-Type", "application/json")

		id := chi.URLParamFromCtx(ctx, "id")
		if id == "" {
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, dto.ErrIDIsEmpty)
			return
		}
		log := h.log.With(slog.String("op", op), slog.String("game_id", id), slog.Int64("user_id", user.ID))

		var req dto.CreateMessageRequest
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, dto.ErrBody)
			return
		}

		exists, err := h.redis.RoomExists(ctx, id)
		if err != nil {
			log.Error("error got room exists", prettylogger.Err(err))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, dto.ErrInternalError)
			return
		}
		if !exists {
			w.WriteHeader(http.StatusNotFound)
			render.JSON(w, r, dto.Error("Игра не найдена"))
			return
		}

		message := &models.Message{
			ID:              uuid.NewString(),
			CreatorID:       user.ID,
			CreatorUsername: user.Username,
			Text:            req.Text,
			CreatedAt:       time.Now().UTC(),
		}
		messageBytes, err := json.Marshal(message)
		if err != nil {
			log.Error("error marshalling message", prettylogger.Err(err))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, dto.ErrInternalError)
			return
		}

		err = h.redis.CreateMessage(ctx, id, messageBytes)
		if err != nil {
			log.Error("error creating message", prettylogger.Err(err))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, dto.ErrInternalError)
			return
		}

		err = h.redis.PublishEvent(ctx, models.Event{
			Type:     models.TypeNewMessage,
			UserID:   user.ID,
			GameID:   id,
			IsPublic: false,
			Payload:  messageBytes,
		})
		if err != nil {
			log.Error("error publishing event", prettylogger.Err(err))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, dto.ErrInternalError)
			return
		}

		log.Info("message created successfully")
		render.JSON(w, r, dto.OK())
	}
}

func (h *Handlers) GetMessages() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.CreateMessage"

		ctx := r.Context()
		user := ctx.Value(middlewares.UserContextKey).(*middlewares.User)
		w.Header().Set("Content-Type", "application/json")

		id := chi.URLParamFromCtx(ctx, "id")
		if id == "" {
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, dto.ErrIDIsEmpty)
			return
		}
		log := h.log.With(slog.String("op", op), slog.String("game_id", id), slog.Int64("user_id", user.ID))

		exists, err := h.redis.RoomExists(ctx, id)
		if err != nil {
			log.Error("error got room exists", prettylogger.Err(err))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, dto.ErrInternalError)
			return
		}
		if !exists {
			w.WriteHeader(http.StatusNotFound)
			render.JSON(w, r, dto.Error("Игра не найдена"))
			return
		}

		messages, err := h.redis.ReadMessages(ctx, id)
		if err != nil {
			log.Error("error read messages", prettylogger.Err(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		log.Info("messages read successfully")
		render.JSON(w, r, dto.ReadMessagesResponse{
			Response: dto.OK(),
			Messages: messages,
		})
	}
}
