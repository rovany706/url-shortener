package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/rovany706/url-shortener/internal/app"
	"github.com/rovany706/url-shortener/internal/config"
	"github.com/rovany706/url-shortener/internal/handlers"
	"github.com/rovany706/url-shortener/internal/middleware"
	"github.com/rovany706/url-shortener/internal/repository"
	"go.uber.org/zap"
)

func MainRouter(app app.URLShortener, appConfig *config.AppConfig, repository repository.Repository, logger *zap.Logger) chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.ResponseLogger(logger))
	r.Use(middleware.RequestLogger(logger))
	r.Use(middleware.RequestGzipCompress())
	r.Use(middleware.ResponseGzipCompress())
	r.Route("/", func(r chi.Router) {
		r.Get("/{id}", handlers.RedirectHandler(app))
		r.Post("/", handlers.MakeShortURLHandler(app, appConfig))
		r.Get("/ping", handlers.PingHandler(repository, logger))
		r.Route("/api", func(r chi.Router) {
			r.Post("/shorten", handlers.MakeShortURLHandlerJSON(app, appConfig, logger))
			r.Post("/shorten/batch", handlers.MakeShortURLBatchHandler(app, appConfig, logger))
		})
	})

	return r
}
