package handlers

import (
	"errors"
	"ms4me/game/internal/http/dto/response"
	userdto "ms4me/game/internal/http/dto/user"
	"ms4me/game/internal/http/middlewares"
	"ms4me/game/internal/services/auth"
	"ms4me/game/internal/storage"
	"ms4me/game/pkg/lib/validator"
	"net/http"

	"github.com/go-chi/render"
)

func (gh *GameHandlers) Register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		w.Header().Add("Content-Type", "application/json")
		var dto userdto.RegisterRequest
		if err := render.DecodeJSON(r.Body, &dto); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, ErrInvalidBody)
			return
		}

		if err := validator.Validate(dto); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Error(validator.GetDetailedError(err).Error()))
			return
		}

		id, err := gh.authSrv.Register(ctx, dto.Username, dto.Password)
		if err != nil {
			if errors.Is(err, storage.ErrUserExists) {
				w.WriteHeader(http.StatusBadRequest)
				render.JSON(w, r, response.Error(storage.ErrUserExists.Error()))
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, response.ErrInternalError)
			return
		}

		render.JSON(w, r, userdto.RegisterResponse{
			Response: response.OK(),
			ID:       id,
		})
	}
}

func (gh *GameHandlers) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		w.Header().Add("Content-Type", "application/json")
		var dto userdto.RegisterRequest
		if err := render.DecodeJSON(r.Body, &dto); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, ErrInvalidBody)
			return
		}

		token, err := gh.authSrv.Login(ctx, dto.Username, dto.Password)
		if err != nil {
			if errors.Is(err, storage.ErrUserNotFound) {
				w.WriteHeader(http.StatusNotFound)
				render.JSON(w, r, ErrUserNotFound)
				return
			}
			if errors.Is(err, auth.ErrIncorrectPassword) {
				w.WriteHeader(http.StatusBadRequest)
				render.JSON(w, r, ErrIncorrectPassword)
			}
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, response.ErrInternalError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "token",
			Value:    token,
			Path:     "/",
			HttpOnly: false,
			Secure:   false,
			SameSite: http.SameSiteStrictMode,
			MaxAge:   3600 * 24 * 30,
		})

		render.JSON(w, r, response.OK())
	}
}

func (gh *GameHandlers) Logout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		cookie := &http.Cookie{
			Name:     "token",
			Value:    "",
			Path:     "/",
			MaxAge:   -1,
			HttpOnly: true,
		}
		http.SetCookie(w, cookie)
		render.JSON(w, r, response.OK())
	}
}

func (h *GameHandlers) User() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		user := ctx.Value(middlewares.UserContextKey).(*middlewares.User)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, userdto.UserResponse{
			Response: response.OK(),
			User: &middlewares.User{
				ID:       user.ID,
				Username: user.Username,
			},
		})
	}
}
