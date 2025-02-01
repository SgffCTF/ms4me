package gamehandlers

import (
	"context"
	"errors"
	gamedto "game-creator/internal/http/dto/game"
	"game-creator/internal/http/dto/response"
	"game-creator/internal/http/dto/validator"
	"game-creator/internal/http/middlewares"
	"game-creator/internal/models"
	"game-creator/internal/storage"
	"log/slog"
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
}

type GameRouter struct {
	log     *slog.Logger
	service GameService
}

func New(log *slog.Logger, service GameService) *GameRouter {
	return &GameRouter{log: log, service: service}
}

func (gr *GameRouter) CreateGame() http.HandlerFunc {
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

func (gr *GameRouter) GetGames() http.HandlerFunc {
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

func (gr *GameRouter) GetGame() http.HandlerFunc {
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
		log = log.With(slog.String("game_id", id))

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

		log.Info("game get successfully", slog.Int64("user_id", user.ID))

		render.JSON(w, r, gamedto.GetGamesResponse{
			Response: response.OK(),
			Games:    []*models.Game{game},
		})
	}
}

func (gr *GameRouter) UpdateGame() http.HandlerFunc {
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

		log.Info("game update successfully", slog.Int64("user_id", user.ID), slog.String("game_id", id))

		render.JSON(w, r, response.OK())
	}
}

func (gr *GameRouter) DeleteGame() http.HandlerFunc {
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
		log = log.With(slog.String("game_id", id))

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

		log.Info("game delete successfully", slog.Int64("user_id", user.ID))

		render.JSON(w, r, response.OK())
	}
}

func (gr *GameRouter) StartGame() http.HandlerFunc {
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
			log.Error("start game error", prettylogger.Err(err), slog.String("id", id))
			render.JSON(w, r, response.Error(response.ErrInternalError.Error()))
			return
		}

		render.JSON(w, r, response.OK())
	}
}

func (gr *GameRouter) EnterGame() http.HandlerFunc {
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
				log.Warn("start game error", prettylogger.Err(err))
				render.JSON(w, r, response.Error(storage.ErrMaxPlayers.Error()))
				return
			}
			if errors.Is(err, storage.ErrPlayerAlreadyExists) {
				log.Warn("start game error", prettylogger.Err(err))
				render.JSON(w, r, response.Error(storage.ErrPlayerAlreadyExists.Error()))
				return
			}
			log.Error("start game error", prettylogger.Err(err), slog.String("id", id))
			render.JSON(w, r, response.Error(response.ErrInternalError.Error()))
			return
		}

		render.JSON(w, r, response.OK())
	}
}
