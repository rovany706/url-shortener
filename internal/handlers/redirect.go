package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rovany706/url-shortener/internal/app"
)

func RedirectHandler(app app.URLShortener) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		shortID := chi.URLParam(r, "id")
		shortenedURLInfo, ok := app.GetFullURL(r.Context(), shortID)
		if ok {
			if shortenedURLInfo.IsDeleted {
				w.WriteHeader(http.StatusGone)
				return
			}
			http.Redirect(w, r, shortenedURLInfo.FullURL, http.StatusTemporaryRedirect)
		} else {
			http.Error(w, "400 Bad Request", http.StatusBadRequest)
		}
	}
}
