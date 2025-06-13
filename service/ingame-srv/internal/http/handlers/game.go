package handlers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"ms4me/game_socket/internal/http/dto"
	"ms4me/game_socket/internal/http/middlewares"
	"ms4me/game_socket/internal/models"
	"ms4me/game_socket/internal/service/game"
	"ms4me/game_socket/pkg/lib/validator"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/jacute/prettylogger"
)

var (
	ErrNotYourGame = dto.Error("Пользователь отсутсвует среди участников игры")
)

func (h *Handlers) GetGameInfo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.GetParticipants"

		ctx := r.Context()
		user := ctx.Value(middlewares.UserContextKey).(*middlewares.User)
		w.Header().Set("Content-Type", "application/json")

		id := chi.URLParamFromCtx(ctx, "id")
		if id == "" {
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, dto.ErrIDIsEmpty)
			return
		}
		log := h.log.With(slog.String("game_id", id), slog.String("op", op), slog.Int64("user_id", user.ID))

		status, err := h.gameClient.GetStatus(id)
		if err != nil {
			log.Error("error getting started info from game-srv", prettylogger.Err(err))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, dto.ErrInternalError)
			return
		}
		if status == "closed" {
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, dto.Error("Игра завершена"))
			return
		}

		roomParticipantsMap, err := h.redis.GetClientsInChannel(ctx, id)
		if err != nil {
			log.Error("error got room participants", prettylogger.Err(err))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, dto.ErrInternalError)
			return
		}

		if _, ok := roomParticipantsMap[strconv.Itoa(int(user.ID))]; !ok {
			log.Info("user not in game")
			w.WriteHeader(http.StatusForbidden)
			render.JSON(w, r, ErrNotYourGame)
			return
		}

		data, err := marshalGameData(roomParticipantsMap)
		if err != nil {
			log.Error("error marshalling room participants", prettylogger.Err(err))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, dto.ErrInternalError)
			return
		}

		render.JSON(w, r, dto.GetParticipantsResponse{
			Response:     dto.OK(),
			Participants: data,
		})
	}
}

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

		userInGame := false
		var loseEvent *models.LoseEvent
		var winEvent *models.WinEvent
		for _, participant := range participants {
			if participant.Field == nil {
				participant.Field = game.CreateField(req.Row, req.Col)
			}

			if participant.ID == user.ID {
				userInGame = true
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

		if !userInGame {
			log.Info("user not in game")
			w.WriteHeader(http.StatusForbidden)
			render.JSON(w, r, ErrNotYourGame)
			return
		}

		payload, err := marshalGameData(participants)
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
			winner := getParticipantWithoutOpenMine(participants)
			if winner == nil {
				log.Error("no winner in room", prettylogger.Err(err))
				w.WriteHeader(http.StatusInternalServerError)
				render.JSON(w, r, dto.ErrInternalError)
				return
			}
			err = h.gameClient.Close(id, winner.ID)
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
			err = h.gameClient.Close(id, winEvent.WinnerID)
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

		userParticipant, ok := participants[strconv.Itoa(int(user.ID))]
		if ok {
			err := userParticipant.Field.SetFlag(req.Row, req.Col)
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
			err = h.redis.AddClientToChannel(ctx, id, userParticipant.ID, userParticipant)
			if err != nil {
				log.Error("error saving participant info", prettylogger.Err(err))
			}
		} else {
			log.Info("user not in game")
			w.WriteHeader(http.StatusForbidden)
			render.JSON(w, r, ErrNotYourGame)
			return
		}

		payload, err := marshalGameData(participants)
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

// marshalGameData подготаливает json с данными по игре для отправки клиенту, маскируя поля json, которые не должны передаваться (расположения мин)
func marshalGameData(participants map[string]*models.RoomParticipant) ([]byte, error) {
	arrParticipants := make([]*models.RoomParticipant, 0)
	for _, participant := range participants {
		if participant.Field != nil {
			for _, row := range participant.Field.Grid {
				for _, c := range row {
					c.NeighborMines = 0
					c.HasMine = nil
				}
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

func getParticipantWithoutOpenMine(participants map[string]*models.RoomParticipant) *models.RoomParticipant {
	for _, rp := range participants {
		if rp.Field != nil && !rp.Field.MineIsOpen {
			return rp
		}
	}
	return nil
}
