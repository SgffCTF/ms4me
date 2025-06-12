package handlers

import (
	gamedto "ms4me/game/internal/http/dto/game"
	"ms4me/game/internal/http/dto/response"
	"ms4me/game/internal/http/middlewares"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

func (h *GameHandlers) OpenCell() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		user := ctx.Value(middlewares.UserContextKey).(*middlewares.User)

		id := chi.URLParamFromCtx(ctx, "id")
		if id == "" {
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, ErrEmptyID)
			return
		}

		var req gamedto.OpenCellRequest
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, ErrInvalidBody)
			return
		}

		if err := h.gameSrv.OpenCell(ctx, &req, id, user.ID); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, response.ErrInternalError)
			return
		}
	}
}
