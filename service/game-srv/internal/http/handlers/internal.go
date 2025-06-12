package handlers

import (
	"errors"
	gamedto "ms4me/game/internal/http/dto/game"
	"ms4me/game/internal/http/dto/response"
	"ms4me/game/internal/storage"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

func (gh *GameHandlers) GameStarted() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		w.Header().Add("Content-Type", "application/json")

		id := chi.URLParam(r, "id")
		if id == "" {
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, ErrEmptyID)
			return
		}

		started, err := gh.gameSrv.GameStarted(ctx, id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, response.ErrInternalError)
			return
		}

		render.JSON(w, r, gamedto.GameStartedResponse{
			Response: response.OK(),
			Started:  started,
		})
	}
}

func (gh *GameHandlers) CloseGame() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		w.Header().Add("Content-Type", "application/json")

		id := chi.URLParam(r, "id")
		if id == "" {
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, ErrEmptyID)
			return
		}

		err := gh.gameSrv.CloseGame(ctx, id)
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

		render.JSON(w, r, response.OK())
	}
}
