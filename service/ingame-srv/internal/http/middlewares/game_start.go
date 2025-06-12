package middlewares

import (
	"log/slog"
	"ms4me/game_socket/internal/http/dto"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/jacute/prettylogger"
)

func (mw *Middlewares) CheckGameStarted() func(next http.Handler) http.Handler {
	const op = "middlewares.CheckGameStarted"

	return func(next http.Handler) http.Handler {
		log := mw.log.With(
			slog.String("op", op),
		)
		log.Info("auth middleware enabled")

		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			id := chi.URLParamFromCtx(ctx, "id")
			if id == "" {
				w.WriteHeader(http.StatusBadRequest)
				render.JSON(w, r, dto.Error("id пустой"))
				return
			}
			started, err := mw.gameClient.Started(id)
			if err != nil {
				log.Error("error getting started info from game-srv", prettylogger.Err(err))
				w.WriteHeader(http.StatusInternalServerError)
				render.JSON(w, r, dto.ErrInternalError)
				return
			}
			if !started {
				w.WriteHeader(http.StatusTooEarly)
				render.JSON(w, r, dto.Error("игра ещё не началась"))
				return
			}

			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}
