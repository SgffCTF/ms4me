package middlewares

import (
	"context"
	"game-creator/internal/http/dto/response"
	"log/slog"
	"net/http"

	ssov1 "game-creator/pkg/grpc/sso"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/jacute/prettylogger"
)

type ContextKey string

var UserContextKey ContextKey = "user"

type User struct {
	ID       int64
	Username string
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
			token := r.Header.Get("Authorization")

			if token == "" {
				render.JSON(w, r, response.Error("Authorization header is empty"))
				return
			}

			data, err := mw.authClient.VerifyToken(r.Context(), &ssov1.VerifyTokenRequest{Token: token})
			if err != nil {
				log.Warn("unauthorized request", prettylogger.Err(err))
				render.JSON(w, r, response.Error("Invalid token"))
				return
			}

			ctx := context.WithValue(r.Context(), UserContextKey, &User{data.User.GetId(), data.User.GetUsername()})
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}
