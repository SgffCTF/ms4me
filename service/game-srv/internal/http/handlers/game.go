package handlers

import (
	"errors"
	gamedto "ms4me/game/internal/http/dto/game"
	"ms4me/game/internal/http/dto/response"
	"ms4me/game/internal/http/middlewares"
	"ms4me/game/internal/services/game"
	"ms4me/game/internal/storage"
	ingameclient "ms4me/game/pkg/ingame_client"
	"ms4me/game/pkg/lib/validator"
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
			if errors.Is(err, storage.ErrAlreadyCreatedGame) {
				render.JSON(w, r, response.Error(storage.ErrAlreadyCreatedGame.Error()))
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

		id := chi.URLParam(r, "id")
		if id == "" {
			render.JSON(w, r, response.Error("id shouldn't be empty"))
			return
		}

		game, err := gr.gameSrv.GetGame(ctx, id)
		if err != nil {
			if errors.Is(err, storage.ErrGameNotFound) {
				w.WriteHeader(http.StatusBadRequest)
				render.JSON(w, r, response.Error(storage.ErrGameNotFound.Error()))
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, response.ErrInternalError)
			return
		}

		render.JSON(w, r, gamedto.GetGameResponse{
			Response: response.OK(),
			Game:     game,
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
				w.WriteHeader(http.StatusBadRequest)
				render.JSON(w, r, response.Error(storage.ErrGameNotFoundOrNotYourOwn.Error()))
				return
			}
			if errors.Is(err, storage.ErrDeleteClosedGame) {
				w.WriteHeader(http.StatusBadRequest)
				render.JSON(w, r, response.Error(storage.ErrDeleteClosedGame.Error()))
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
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
			w.WriteHeader(http.StatusBadRequest)
			if errors.Is(err, game.ErrOnlyOwnerCanStartGame) {
				render.JSON(w, r, response.Error(game.ErrOnlyOwnerCanStartGame.Error()))
				return
			}
			if errors.Is(err, game.ErrGameIsNotOpen) {
				render.JSON(w, r, response.Error(game.ErrGameIsNotOpen.Error()))
				return
			}
			if errors.Is(err, game.ErrGameAlreadyStarted) {
				render.JSON(w, r, response.Error(game.ErrGameAlreadyStarted.Error()))
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
			if errors.Is(err, ingameclient.ErrNotReady) {
				render.JSON(w, r, response.Error(ingameclient.ErrNotReady.Error()))
				return
			}
			if errors.Is(err, storage.ErrGameNotFoundOrNotYourOwn) {
				render.JSON(w, r, response.Error(storage.ErrGameNotFound.Error()))
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
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

		err := gr.gameSrv.EnterGame(ctx, id, user.ID, user.Username)
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
			if errors.Is(err, game.ErrGameIsNotOpen) {
				render.JSON(w, r, response.Error(game.ErrGameIsNotOpen.Error()))
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

		err := gr.gameSrv.ExitGame(ctx, id, user.ID, user.Username)
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

func (gr *GameHandlers) GetMyGames() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		user := ctx.Value(middlewares.UserContextKey).(*middlewares.User)
		w.Header().Add("Content-Type", "application/json")

		games, err := gr.gameSrv.UserGames(ctx, user.ID)
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

func (gr *GameHandlers) GetCongratulation() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		w.Header().Add("Content-Type", "application/json")

		id := chi.URLParam(r, "id")
		if id == "" {
			render.JSON(w, r, ErrEmptyID)
			return
		}

		data, err := gr.gameSrv.Congratulation(ctx, id)
		if err != nil {
			if errors.Is(err, game.ErrGameIsNotClosed) {
				w.WriteHeader(http.StatusBadRequest)
				render.JSON(w, r, response.Error(game.ErrGameIsNotClosed.Error()))
				return
			}
			if errors.Is(err, game.ErrTemplate) {
				w.WriteHeader(http.StatusBadRequest)
				render.JSON(w, r, response.Error(game.ErrTemplate.Error()))
				return
			}
			if errors.Is(err, storage.ErrGameNotFound) {
				w.WriteHeader(http.StatusNotFound)
				render.JSON(w, r, response.Error(storage.ErrGameNotFound.Error()))
			}
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, response.ErrInternalError)
			return
		}

		_, err = w.Write(data)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, response.Error("Ошибка при ответе"))
			return
		}
	}
}
