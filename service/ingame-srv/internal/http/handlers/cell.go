package handlers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"ms4me/game_socket/internal/game"
	"ms4me/game_socket/internal/http/dto"
	"ms4me/game_socket/internal/http/middlewares"
	"ms4me/game_socket/internal/models"
	"ms4me/game_socket/pkg/lib/validator"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/jacute/prettylogger"
)

func (h *Handlers) OpenCell() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.OpenCell"

		ctx := r.Context()
		user := ctx.Value(middlewares.UserContextKey).(*middlewares.User)
		log := h.log.With(slog.String("op", op), slog.Int64("user_id", user.ID))

		var req dto.ClickCellRequest
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, dto.ErrBody)
			return
		}

		if err := validator.Validate(req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, dto.Error(validator.GetDetailedError(err).Error()))
			return
		}

		id := chi.URLParamFromCtx(ctx, "id")
		log = log.With(slog.String("game_id", id))
		participants, err := h.redis.GetClientsInChannel(ctx, id)
		if err != nil {
			log.Error("error getting room participants", prettylogger.Err(err))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, dto.ErrInternalError)
			return
		}

		var loseEvent *models.LoseEvent
		var winEvent *models.WinEvent
		for _, participant := range participants {
			if participant.Field == nil {
				participant.Field = game.CreateField(req.Row, req.Col)
			}

			if participant.ID == user.ID {
				err := participant.Field.OpenCell(req.Row, req.Col)
				if err != nil {
					if errors.Is(err, game.ErrAlreadyOpen) {
						log.Debug("cell already open")
						w.WriteHeader(http.StatusBadRequest)
						render.JSON(w, r, dto.Error(game.ErrAlreadyOpen.Error()))
						return
					}
					if errors.Is(err, game.ErrFieldSize) {
						log.Debug("open outside the field")
						render.JSON(w, r, dto.Error(game.ErrFieldSize.Error()))
						return
					}
					log.Error("error opening cell", prettylogger.Err(err))
					w.WriteHeader(http.StatusInternalServerError)
					render.JSON(w, r, dto.ErrInternalError)
					return
				}
				if participant.Field.MineIsOpen {
					loseEvent = &models.LoseEvent{
						LoserID:       participant.ID,
						LoserUsername: participant.Username,
					}
				} else if participant.Field.IsWin() {
					winEvent = &models.WinEvent{
						WinnerID:       participant.ID,
						WinnerUsername: participant.Username,
					}
				}
			}

			err := h.redis.AddClientToChannel(ctx, id, participant.ID, participant)
			if err != nil {
				log.Error("error saving participant info", prettylogger.Err(err))
				w.WriteHeader(http.StatusInternalServerError)
				render.JSON(w, r, dto.ErrInternalError)
				return
			}
		}

		payload, err := marshalClickEventPayload(participants)
		if err != nil {
			log.Error("error marshalling field", prettylogger.Err(err))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, dto.ErrInternalError)
			return
		}
		err = h.redis.PublishEvent(ctx, models.Event{
			Type:     models.TypeClickGame,
			UserID:   user.ID,
			GameID:   id,
			IsPublic: false,
			Payload:  payload,
		})
		if err != nil {
			log.Error("error publishing event", prettylogger.Err(err))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, dto.ErrInternalError)
			return
		}

		if loseEvent != nil {
			resultMarshalled, err := json.Marshal(&loseEvent)
			if err != nil {
				log.Error("error marshalling result", prettylogger.Err(err))
				w.WriteHeader(http.StatusInternalServerError)
				render.JSON(w, r, dto.ErrInternalError)
				return
			}
			err = h.redis.PublishEvent(ctx, models.Event{
				Type:     models.TypeLoseGame,
				UserID:   user.ID,
				GameID:   id,
				IsPublic: false,
				Payload:  resultMarshalled,
			})
			if err != nil {
				log.Error("error publishing event", prettylogger.Err(err))
				w.WriteHeader(http.StatusInternalServerError)
				render.JSON(w, r, dto.ErrInternalError)
				return
			}
			err = h.gameClient.Close(id)
			if err != nil {
				log.Error("error closing game", prettylogger.Err(err))
				w.WriteHeader(http.StatusInternalServerError)
				render.JSON(w, r, dto.ErrInternalError)
				return
			}
		}
		if winEvent != nil {
			resultMarshalled, err := json.Marshal(&winEvent)
			if err != nil {
				log.Error("error marshalling result", prettylogger.Err(err))
				w.WriteHeader(http.StatusInternalServerError)
				render.JSON(w, r, dto.ErrInternalError)
				return
			}
			err = h.redis.PublishEvent(ctx, models.Event{
				Type:     models.TypeWinGame,
				UserID:   user.ID,
				GameID:   id,
				IsPublic: false,
				Payload:  resultMarshalled,
			})
			if err != nil {
				log.Error("error publishing event", prettylogger.Err(err))
				w.WriteHeader(http.StatusInternalServerError)
				render.JSON(w, r, dto.ErrInternalError)
				return
			}
			err = h.gameClient.Close(id)
			if err != nil {
				log.Error("error closing game", prettylogger.Err(err))
				w.WriteHeader(http.StatusInternalServerError)
				render.JSON(w, r, dto.ErrInternalError)
				return
			}
		}

		render.JSON(w, r, dto.OK())
	}
}

func (h *Handlers) Flag() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.Flag"
		ctx := r.Context()
		user := ctx.Value(middlewares.UserContextKey).(*middlewares.User)
		log := h.log.With(slog.String("op", op), slog.Int64("user_id", user.ID))

		var req dto.ClickCellRequest
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, dto.ErrBody)
			return
		}

		if err := validator.Validate(req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, dto.Error(validator.GetDetailedError(err).Error()))
			return
		}

		id := chi.URLParamFromCtx(ctx, "id")
		log = log.With(slog.String("game_id", id))
		participants, err := h.redis.GetClientsInChannel(ctx, id)
		if err != nil {
			log.Error("error getting room participants", prettylogger.Err(err))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, dto.ErrInternalError)
			return
		}

		for _, participant := range participants {
			if participant.ID == user.ID {
				err := participant.Field.SetFlag(req.Row, req.Col)
				if err != nil {
					if errors.Is(err, game.ErrFlagOnOpenCell) {
						log.Debug("error setting flag on open cell")
						w.WriteHeader(http.StatusBadRequest)
						render.JSON(w, r, dto.Error(game.ErrFlagOnOpenCell.Error()))
						return
					}
					if errors.Is(err, game.ErrFieldSize) {
						log.Debug("open outside the field")
						render.JSON(w, r, dto.Error(game.ErrFieldSize.Error()))
						return
					}
					log.Error("error setting flag", prettylogger.Err(err))
					w.WriteHeader(http.StatusInternalServerError)
					render.JSON(w, r, dto.ErrInternalError)
					return
				}
				err = h.redis.AddClientToChannel(ctx, id, participant.ID, participant)
				if err != nil {
					log.Error("error saving participant info", prettylogger.Err(err))
				}
			}
		}
		payload, err := marshalClickEventPayload(participants)
		if err != nil {
			log.Error("error marshalling field", prettylogger.Err(err))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, dto.ErrInternalError)
			return
		}
		err = h.redis.PublishEvent(ctx, models.Event{
			Type:     models.TypeClickGame,
			UserID:   user.ID,
			GameID:   id,
			IsPublic: false,
			Payload:  payload,
		})
		if err != nil {
			h.log.Error("error publishing event", prettylogger.Err(err))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, dto.ErrInternalError)
			return
		}
		render.JSON(w, r, dto.OK())
	}
}

// marshalClickEventPayload подготаливает event для пуша в очередь redis, маскируя поля json, которые не должны передаваться по ws
func marshalClickEventPayload(participants map[string]*models.RoomParticipant) ([]byte, error) {
	arrParticipants := make([]*models.RoomParticipant, 0)
	for _, participant := range participants {
		for _, row := range participant.Field.Grid {
			for _, c := range row {
				c.NeighborMines = 0
				c.HasMine = nil
			}
		}
		arrParticipants = append(arrParticipants, participant)
	}
	data, err := json.Marshal(arrParticipants)
	if err != nil {
		return nil, err
	}
	return data, nil
}
