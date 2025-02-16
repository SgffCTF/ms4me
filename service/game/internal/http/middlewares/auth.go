package middlewares

import (
	"context"
	ssov1 "ms4me/game/pkg/grpc/sso"
	"net/http"
)

type User struct {
	ID       int64
	Username string
}

type contextKey string

const userContextKey contextKey = "user"

func (mw *Middlewares) Auth() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			token := r.Header.Get("Authorization")

			if token == "" {
				http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
				return
			}

			response, err := mw.authSrv.VerifyToken(ctx, &ssov1.VerifyTokenRequest{Token: token})
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			user := response.GetUser()
			ctx = context.WithValue(ctx, userContextKey, &User{ID: user.Id, Username: user.Username})
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
