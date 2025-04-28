package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rovany706/url-shortener/internal/app"
)

// RedirectHandlers обработчики методов перенаправления
type RedirectHandlers struct {
	app app.URLShortener
}

// NewRedirectHandlers создает RedirectHandlers
func NewRedirectHandlers(app app.URLShortener) RedirectHandlers {
	return RedirectHandlers{
		app: app,
	}
}

// RedirectHandler хэндлер перенаправления сокращенной ссылки
func (h *RedirectHandlers) RedirectHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		shortID := chi.URLParam(r, "id")
		shortenedURLInfo, ok := h.app.GetFullURL(r.Context(), shortID)
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
