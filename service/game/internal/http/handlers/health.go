package handlers

import (
	"net/http"

	"github.com/go-chi/render"
)

func Health() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		render.Data(w, r, []byte("OK"))
	}
}
