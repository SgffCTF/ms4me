package middlewares

import (
	"context"
	"log/slog"
	"ms4me/game/internal/http/dto/response"
	"ms4me/game/internal/lib/jwt"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/jacute/prettylogger"
)

var (
	ErrInvalidToken = response.Error("Токен не валиден")
	ErrEmptyToken   = response.Error("Токен отсутствует")
)

type ContextKey string

var UserContextKey ContextKey = "user"

type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
}

func (mw *Middlewares) Auth() func(next http.Handler) http.Handler {
	const op = "middlewares.Auth"

	return func(next http.Handler) http.Handler {
		log := mw.log.With(
			slog.String("op", op),
		)
		log.Info("auth middleware enabled")

		fn := func(w http.ResponseWriter, r *http.Request) {
			rid := middleware.GetReqID(r.Context())
			log = log.With(slog.String("rid", rid))
			tokenCookie := r.CookiesNamed("token")
			if len(tokenCookie) == 0 {
				render.JSON(w, r, ErrEmptyToken)
				return
			}
			token := tokenCookie[0].Value

			data, err := jwt.VerifyToken(token, mw.jwtSecret)
			if err != nil {
				log.Warn("unauthorized request", prettylogger.Err(err))
				render.JSON(w, r, ErrInvalidToken)
				return
			}
			id, ok1 := data["user_id"].(float64)
			username, ok2 := data["username"].(string)
			if !ok1 || !ok2 {
				log.Warn("unauthorized request")
				render.JSON(w, r, ErrInvalidToken)
				return
			}

			ctx := context.WithValue(r.Context(), UserContextKey, &User{ID: int64(id), Username: username})
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}
