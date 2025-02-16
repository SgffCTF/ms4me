package gamehandlers

import (
	"context"
	"errors"
	"log/slog"
	gamedto "ms4me/game_creator/internal/http/dto/game"
	"ms4me/game_creator/internal/http/dto/response"
	"ms4me/game_creator/internal/http/dto/validator"
	"ms4me/game_creator/internal/http/middlewares"
	"ms4me/game_creator/internal/models"
	"ms4me/game_creator/internal/storage"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/jacute/prettylogger"
)

var (
	ErrEmptyID = errors.New("id shouldn't be empty")
)

type GameService interface {
	CreateGame(ctx context.Context, userID int64, game *gamedto.CreateGameRequest) (string, error)
	GetGames(ctx context.Context, filter *gamedto.GetGamesRequest) ([]*models.Game, error)
	GetGame(ctx context.Context, id string, userID int64) (*models.Game, error)
	UpdateGame(ctx context.Context, id string, userID int64, game *gamedto.UpdateGameRequest) error
	DeleteGame(ctx context.Context, id string, userID int64) error
	StartGame(ctx context.Context, id string, userID int64) error
	EnterGame(ctx context.Context, id string, userID int64) error
	ExitGame(ctx context.Context, id string, userID int64) error
}

type GameHandlers struct {
	log     *slog.Logger
	service GameService
}

func New(log *slog.Logger, service GameService) *GameHandlers {
	return &GameHandlers{log: log, service: service}
}

func (gr *GameHandlers) CreateGame() http.HandlerFunc {
	const op = "handlers.CreateGame"
	log := gr.log.With(slog.String("op", op))
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		user := ctx.Value(middlewares.UserContextKey).(*middlewares.User)
		log = log.With(slog.Int64("user_id", user.ID))

		var req gamedto.CreateGameRequest
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Warn("invalid create game request", prettylogger.Err(err))
			render.JSON(w, r, response.Error("invalid request"))
			return
		}

		if err := req.Validate(); err != nil {
			render.JSON(w, r, response.Error(validator.GetDetailedError(err)))
			return
		}

		id, err := gr.service.CreateGame(ctx, user.ID, &req)
		if err != nil {
			if errors.Is(err, storage.ErrAlreadyPlaying) {
				log.Warn("create game error", prettylogger.Err(err))
				render.JSON(w, r, response.Error(storage.ErrAlreadyPlaying.Error()))
				return
			}
			log.Error("create game error", prettylogger.Err(err))
			render.JSON(w, r, response.Error(response.ErrInternalError.Error()))
			return
		}

		log.Info(
			"game created successfully",
			slog.Int64("user_id", user.ID),
			slog.String("game_id", id),
		)

		render.JSON(w, r, gamedto.CreateGameResponse{
			Response: response.OK(),
			ID:       id,
		})
	}
}

func (gr *GameHandlers) GetGames() http.HandlerFunc {
	const op = "handlers.CreateGame"
	log := gr.log.With(slog.String("op", op))
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		user := ctx.Value(middlewares.UserContextKey).(*middlewares.User)

		var req gamedto.GetGamesRequest
		if err := req.Render(r.URL.Query()); err != nil {
			render.JSON(w, r, response.Error(err.Error()))
			return
		}

		games, err := gr.service.GetGames(ctx, &req)
		if err != nil {
			log.Error("get games error", prettylogger.Err(err), slog.Any("dto", req))
			render.JSON(w, r, response.Error(response.ErrInternalError.Error()))
			return
		}

		log.Info("games get successfully", slog.Int64("user_id", user.ID))

		render.JSON(w, r, gamedto.GetGamesResponse{
			Response: response.OK(),
			Games:    games,
		})
	}
}

func (gr *GameHandlers) GetGame() http.HandlerFunc {
	const op = "handlers.GetGame"
	log := gr.log.With(slog.String("op", op))
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		user := ctx.Value(middlewares.UserContextKey).(*middlewares.User)

		id := chi.URLParam(r, "id")
		if id == "" {
			render.JSON(w, r, response.Error("id shouldn't be empty"))
			return
		}
		log = log.With(slog.String("game_id", id), slog.Int64("user_id", user.ID))

		game, err := gr.service.GetGame(ctx, id, user.ID)
		if err != nil {
			if errors.Is(err, storage.ErrGameNotFoundOrNotYourOwn) {
				log.Warn("game not found or not owner")
				render.JSON(w, r, response.Error(storage.ErrGameNotFoundOrNotYourOwn.Error()))
				return
			}
			log.Error("get game error", prettylogger.Err(err))
			render.JSON(w, r, response.Error(response.ErrInternalError.Error()))
			return
		}

		log.Info("game get successfully")

		render.JSON(w, r, gamedto.GetGamesResponse{
			Response: response.OK(),
			Games:    []*models.Game{game},
		})
	}
}

func (gr *GameHandlers) UpdateGame() http.HandlerFunc {
	const op = "handlers.UpdateGame"
	log := gr.log.With(slog.String("op", op))
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		user := ctx.Value(middlewares.UserContextKey).(*middlewares.User)

		id := chi.URLParam(r, "id")
		if id == "" {
			render.JSON(w, r, response.Error("id shouldn't be empty"))
			return
		}
		log = log.With(slog.Int64("user_id", user.ID), slog.String("game_id", id))

		var req gamedto.UpdateGameRequest
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			render.JSON(w, r, response.Error(err.Error()))
			return
		}

		if err := req.Validate(); err != nil {
			log.Warn("invalid update game request", prettylogger.Err(err))
			render.JSON(w, r, response.Error(validator.GetDetailedError(err)))
			return
		}

		err := gr.service.UpdateGame(ctx, id, user.ID, &req)
		if err != nil {
			if errors.Is(err, storage.ErrGameNotFoundOrNotYourOwn) {
				log.Warn("update game error", prettylogger.Err(err))
				render.JSON(w, r, response.Error(storage.ErrGameNotFoundOrNotYourOwn.Error()))
				return
			}
			log.Error("update game error", prettylogger.Err(err), slog.String("id", id))
			render.JSON(w, r, response.Error(response.ErrInternalError.Error()))
			return
		}

		log.Info("game update successfully")

		render.JSON(w, r, response.OK())
	}
}

func (gr *GameHandlers) DeleteGame() http.HandlerFunc {
	const op = "handlers.DeleteGame"
	log := gr.log.With(slog.String("op", op))
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		user := ctx.Value(middlewares.UserContextKey).(*middlewares.User)

		id := chi.URLParam(r, "id")
		if id == "" {
			render.JSON(w, r, response.Error("id shouldn't be empty"))
			return
		}
		log = log.With(slog.Int64("user_id", user.ID), slog.String("game_id", id))

		err := gr.service.DeleteGame(ctx, id, user.ID)
		if err != nil {
			if errors.Is(err, storage.ErrGameNotFoundOrNotYourOwn) {
				log.Warn("delete game error", prettylogger.Err(err))
				render.JSON(w, r, response.Error(storage.ErrGameNotFoundOrNotYourOwn.Error()))
				return
			}
			log.Error("delete game error", prettylogger.Err(err))
			render.JSON(w, r, response.Error(response.ErrInternalError.Error()))
			return
		}

		log.Info("game delete successfully")

		render.JSON(w, r, response.OK())
	}
}

func (gr *GameHandlers) StartGame() http.HandlerFunc {
	const op = "handlers.StartGame"
	log := gr.log.With(slog.String("op", op))
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		user := ctx.Value(middlewares.UserContextKey).(*middlewares.User)

		id := chi.URLParam(r, "id")
		if id == "" {
			render.JSON(w, r, response.Error(ErrEmptyID.Error()))
			return
		}
		log = log.With(slog.Int64("user_id", user.ID), slog.String("game_id", id))

		err := gr.service.StartGame(ctx, id, user.ID)
		if err != nil {
			if errors.Is(err, storage.ErrOnlyOwnerCanStartGame) {
				log.Warn("start game error", prettylogger.Err(err))
				render.JSON(w, r, response.Error(storage.ErrOnlyOwnerCanStartGame.Error()))
				return
			}
			if errors.Is(err, storage.ErrGameIsNotOpen) {
				log.Warn("start game error", prettylogger.Err(err))
				render.JSON(w, r, response.Error(storage.ErrGameIsNotOpen.Error()))
				return
			}
			if errors.Is(err, storage.ErrGameAlreadyStarted) {
				log.Warn("start game error", prettylogger.Err(err))
				render.JSON(w, r, response.Error(storage.ErrGameAlreadyStarted.Error()))
				return
			}
			if errors.Is(err, storage.ErrIncorrectCountOfPlayers) {
				log.Warn("start game error", prettylogger.Err(err))
				render.JSON(w, r, response.Error(storage.ErrIncorrectCountOfPlayers.Error()))
				return
			}
			if errors.Is(err, storage.ErrGameNotFound) {
				log.Warn("start game error", prettylogger.Err(err))
				render.JSON(w, r, response.Error(storage.ErrGameNotFound.Error()))
				return
			}
			log.Error("start game error", prettylogger.Err(err))
			render.JSON(w, r, response.Error(response.ErrInternalError.Error()))
			return
		}

		log.Info("game started successfully")

		render.JSON(w, r, response.OK())
	}
}

func (gr *GameHandlers) EnterGame() http.HandlerFunc {
	const op = "handlers.EnterGame"
	log := gr.log.With(slog.String("op", op))
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		user := ctx.Value(middlewares.UserContextKey).(*middlewares.User)

		id := chi.URLParam(r, "id")
		if id == "" {
			render.JSON(w, r, response.Error(ErrEmptyID.Error()))
			return
		}
		log = log.With(slog.Int64("user_id", user.ID), slog.String("game_id", id))

		err := gr.service.EnterGame(ctx, id, user.ID)
		if err != nil {
			if errors.Is(err, storage.ErrMaxPlayers) {
				log.Warn("enter game error", prettylogger.Err(err))
				render.JSON(w, r, response.Error(storage.ErrMaxPlayers.Error()))
				return
			}
			if errors.Is(err, storage.ErrPlayerAlreadyExists) {
				log.Warn("enter game error", prettylogger.Err(err))
				render.JSON(w, r, response.Error(storage.ErrPlayerAlreadyExists.Error()))
				return
			}
			if errors.Is(err, storage.ErrAlreadyPlaying) {
				log.Warn("enter game error", prettylogger.Err(err))
				render.JSON(w, r, response.Error(storage.ErrAlreadyPlaying.Error()))
				return
			}
			if errors.Is(err, storage.ErrGameIsNotOpen) {
				log.Warn("enter game error", prettylogger.Err(err))
				render.JSON(w, r, response.Error(storage.ErrGameIsNotOpen.Error()))
				return
			}
			if errors.Is(err, storage.ErrGameNotFound) {
				log.Warn("enter game error", prettylogger.Err(err))
				render.JSON(w, r, response.Error(storage.ErrGameNotFound.Error()))
				return
			}
			log.Error("enter game error", prettylogger.Err(err))
			render.JSON(w, r, response.Error(response.ErrInternalError.Error()))
			return
		}

		log.Info("game entered successfully")

		render.JSON(w, r, response.OK())
	}
}

func (gr *GameHandlers) ExitGame() http.HandlerFunc {
	const op = "handlers.ExitGame"
	log := gr.log.With(slog.String("op", op))
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		user := ctx.Value(middlewares.UserContextKey).(*middlewares.User)

		id := chi.URLParam(r, "id")
		if id == "" {
			render.JSON(w, r, response.Error(ErrEmptyID.Error()))
			return
		}
		log = log.With(slog.Int64("user_id", user.ID), slog.String("game_id", id))

		err := gr.service.ExitGame(ctx, id, user.ID)
		if err != nil {
			if errors.Is(err, storage.ErrGameNotFound) {
				log.Warn("exit game error", prettylogger.Err(err))
				render.JSON(w, r, response.Error(storage.ErrGameNotFound.Error()))
				return
			}
			if errors.Is(err, storage.ErrOwnerCantExitFromOwnGame) {
				log.Warn("exit game error", prettylogger.Err(err))
				render.JSON(w, r, response.Error(storage.ErrOwnerCantExitFromOwnGame.Error()))
				return
			}
			if errors.Is(err, storage.ErrYouNotParticipate) {
				log.Warn("exit game error", prettylogger.Err(err))
				render.JSON(w, r, response.Error(storage.ErrYouNotParticipate.Error()))
				return
			}
			if errors.Is(err, storage.ErrGameNotFound) {
				log.Warn("exit game error", prettylogger.Err(err))
				render.JSON(w, r, response.Error(storage.ErrGameNotFound.Error()))
				return
			}
			log.Error("exit game error", prettylogger.Err(err))
			render.JSON(w, r, response.Error(response.ErrInternalError.Error()))
			return
		}

		log.Info("exited from game successfully")

		render.JSON(w, r, response.OK())
	}
}
