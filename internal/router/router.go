package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/rovany706/url-shortener/internal/app"
	"github.com/rovany706/url-shortener/internal/config"
	"github.com/rovany706/url-shortener/internal/handlers"
	"github.com/rovany706/url-shortener/internal/middleware"
	"go.uber.org/zap"
)

func MainRouter(app app.URLShortener, appConfig *config.AppConfig, logger *zap.Logger) chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.ResponseLogger(logger))
	r.Use(middleware.RequestLogger(logger))
	r.Route("/", func(r chi.Router) {
		r.Get("/{id}", handlers.RedirectHandler(app))
		r.Post("/", handlers.MakeShortURLHandler(app, appConfig))
		r.Post("/api/shorten", handlers.MakeShortURLHandlerJSON(app, appConfig, logger))
	})

	return r
}
