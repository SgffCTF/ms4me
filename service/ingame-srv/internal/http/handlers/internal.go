package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
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
			result := h.wsSrv.CheckRoomConn(participant.ID, id)
			if !result {
				w.WriteHeader(http.StatusTooEarly) // игроки не готовы к началу игры
				return
			}
		}
		w.WriteHeader(http.StatusOK)
	}
}
