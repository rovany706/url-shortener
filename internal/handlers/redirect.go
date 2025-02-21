package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rovany706/url-shortener/internal/app"
)

func RedirectHandler(app app.URLShortener) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		shortID := chi.URLParam(r, "id")
		if fullURL, ok := app.GetFullURL(r.Context(), shortID); ok {
			http.Redirect(w, r, fullURL, http.StatusTemporaryRedirect)
		} else {
			http.Error(w, "400 Bad Request", http.StatusBadRequest)
		}
	}
}
