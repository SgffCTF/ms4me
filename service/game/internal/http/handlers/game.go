package handlers

import (
	"errors"
	gamedto "ms4me/game/internal/http/dto/game"
	"ms4me/game/internal/http/dto/response"
	"ms4me/game/internal/http/middlewares"
	"ms4me/game/internal/lib/validator"
	"ms4me/game/internal/models"
	"ms4me/game/internal/storage"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

func (gr *GameHandlers) CreateGame() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		ctx := r.Context()
		user := ctx.Value(middlewares.UserContextKey).(*middlewares.User)

		var req gamedto.CreateGameRequest
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			render.JSON(w, r, response.Error("invalid request"))
			return
		}

		if err := req.Validate(); err != nil {
			render.JSON(w, r, response.Error(validator.GetDetailedError(err).Error()))
			return
		}

		id, err := gr.gameSrv.CreateGame(ctx, user.ID, &req)
		if err != nil {
			if errors.Is(err, storage.ErrAlreadyPlaying) {
				render.JSON(w, r, response.Error(storage.ErrAlreadyPlaying.Error()))
				return
			}
			render.JSON(w, r, response.ErrInternalError)
			return
		}

		render.JSON(w, r, gamedto.CreateGameResponse{
			Response: response.OK(),
			ID:       id,
		})
	}
}

func (gr *GameHandlers) GetGames() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		ctx := r.Context()

		var req gamedto.GetGamesRequest
		if err := req.Render(r.URL.Query()); err != nil {
			render.JSON(w, r, response.Error(err.Error()))
			return
		}

		games, err := gr.gameSrv.GetGames(ctx, &req)
		if err != nil {
			render.JSON(w, r, response.ErrInternalError)
			return
		}

		render.JSON(w, r, gamedto.GetGamesResponse{
			Response: response.OK(),
			Games:    games,
		})
	}
}

func (gr *GameHandlers) GetGame() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		ctx := r.Context()
		user := ctx.Value(middlewares.UserContextKey).(*middlewares.User)

		id := chi.URLParam(r, "id")
		if id == "" {
			render.JSON(w, r, response.Error("id shouldn't be empty"))
			return
		}

		game, err := gr.gameSrv.GetGame(ctx, id, user.ID)
		if err != nil {
			if errors.Is(err, storage.ErrGameNotFoundOrNotYourOwn) {
				w.WriteHeader(http.StatusBadRequest)
				render.JSON(w, r, response.Error(storage.ErrGameNotFoundOrNotYourOwn.Error()))
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, response.ErrInternalError)
			return
		}

		render.JSON(w, r, gamedto.GetGamesResponse{
			Response: response.OK(),
			Games:    []*models.Game{game},
		})
	}
}

func (gr *GameHandlers) UpdateGame() http.HandlerFunc {
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
			render.JSON(w, r, response.Error(validator.GetDetailedError(err).Error()))
			return
		}

		err := gr.gameSrv.UpdateGame(ctx, id, user.ID, &req)
		if err != nil {
			if errors.Is(err, storage.ErrGameNotFoundOrNotYourOwn) {
				render.JSON(w, r, response.Error(storage.ErrGameNotFoundOrNotYourOwn.Error()))
				return
			}
			render.JSON(w, r, response.ErrInternalError)
			return
		}

		render.JSON(w, r, response.OK())
	}
}

func (gr *GameHandlers) DeleteGame() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		user := ctx.Value(middlewares.UserContextKey).(*middlewares.User)

		id := chi.URLParam(r, "id")
		if id == "" {
			render.JSON(w, r, response.Error("id shouldn't be empty"))
			return
		}

		err := gr.gameSrv.DeleteGame(ctx, id, user.ID)
		if err != nil {
			if errors.Is(err, storage.ErrGameNotFoundOrNotYourOwn) {
				render.JSON(w, r, response.Error(storage.ErrGameNotFoundOrNotYourOwn.Error()))
				return
			}
			render.JSON(w, r, response.ErrInternalError)
			return
		}

		render.JSON(w, r, response.OK())
	}
}

func (gr *GameHandlers) StartGame() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		user := ctx.Value(middlewares.UserContextKey).(*middlewares.User)

		id := chi.URLParam(r, "id")
		if id == "" {
			render.JSON(w, r, ErrEmptyID)
			return
		}

		err := gr.gameSrv.StartGame(ctx, id, user.ID)
		if err != nil {
			if errors.Is(err, storage.ErrOnlyOwnerCanStartGame) {
				render.JSON(w, r, response.Error(storage.ErrOnlyOwnerCanStartGame.Error()))
				return
			}
			if errors.Is(err, storage.ErrGameIsNotOpen) {
				render.JSON(w, r, response.Error(storage.ErrGameIsNotOpen.Error()))
				return
			}
			if errors.Is(err, storage.ErrGameAlreadyStarted) {
				render.JSON(w, r, response.Error(storage.ErrGameAlreadyStarted.Error()))
				return
			}
			if errors.Is(err, storage.ErrIncorrectCountOfPlayers) {
				render.JSON(w, r, response.Error(storage.ErrIncorrectCountOfPlayers.Error()))
				return
			}
			if errors.Is(err, storage.ErrGameNotFound) {
				render.JSON(w, r, response.Error(storage.ErrGameNotFound.Error()))
				return
			}
			render.JSON(w, r, response.ErrInternalError)
			return
		}

		render.JSON(w, r, response.OK())
	}
}

func (gr *GameHandlers) EnterGame() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		user := ctx.Value(middlewares.UserContextKey).(*middlewares.User)

		id := chi.URLParam(r, "id")
		if id == "" {
			render.JSON(w, r, ErrEmptyID)
			return
		}

		err := gr.gameSrv.EnterGame(ctx, id, user.ID)
		if err != nil {
			if errors.Is(err, storage.ErrMaxPlayers) {
				render.JSON(w, r, response.Error(storage.ErrMaxPlayers.Error()))
				return
			}
			if errors.Is(err, storage.ErrPlayerAlreadyExists) {
				render.JSON(w, r, response.Error(storage.ErrPlayerAlreadyExists.Error()))
				return
			}
			if errors.Is(err, storage.ErrAlreadyPlaying) {
				render.JSON(w, r, response.Error(storage.ErrAlreadyPlaying.Error()))
				return
			}
			if errors.Is(err, storage.ErrGameIsNotOpen) {
				render.JSON(w, r, response.Error(storage.ErrGameIsNotOpen.Error()))
				return
			}
			if errors.Is(err, storage.ErrGameNotFound) {
				render.JSON(w, r, response.Error(storage.ErrGameNotFound.Error()))
				return
			}
			render.JSON(w, r, response.ErrInternalError)
			return
		}

		render.JSON(w, r, response.OK())
	}
}

func (gr *GameHandlers) ExitGame() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		user := ctx.Value(middlewares.UserContextKey).(*middlewares.User)

		id := chi.URLParam(r, "id")
		if id == "" {
			render.JSON(w, r, ErrEmptyID)
			return
		}

		err := gr.gameSrv.ExitGame(ctx, id, user.ID)
		if err != nil {
			if errors.Is(err, storage.ErrGameNotFound) {
				render.JSON(w, r, response.Error(storage.ErrGameNotFound.Error()))
				return
			}
			if errors.Is(err, storage.ErrOwnerCantExitFromOwnGame) {
				render.JSON(w, r, response.Error(storage.ErrOwnerCantExitFromOwnGame.Error()))
				return
			}
			if errors.Is(err, storage.ErrYouNotParticipate) {
				render.JSON(w, r, response.Error(storage.ErrYouNotParticipate.Error()))
				return
			}
			if errors.Is(err, storage.ErrGameNotFound) {
				render.JSON(w, r, response.Error(storage.ErrGameNotFound.Error()))
				return
			}
			render.JSON(w, r, response.ErrInternalError)
			return
		}

		render.JSON(w, r, response.OK())
	}
}
