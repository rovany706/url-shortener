package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/rovany706/url-shortener/internal/app"
	"github.com/rovany706/url-shortener/internal/config"
	"github.com/rovany706/url-shortener/internal/handlers"
)

func MainRouter(app app.URLShortener, appConfig *config.AppConfig) chi.Router {
	r := chi.NewRouter()
	r.Route("/", func(r chi.Router) {
		r.Get("/{id}", handlers.RedirectHandler(app))
		r.Post("/", handlers.MakeShortURLHandler(app, appConfig))
	})

	return r
}
