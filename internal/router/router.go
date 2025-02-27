package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/rovany706/url-shortener/internal/app"
	"github.com/rovany706/url-shortener/internal/auth"
	"github.com/rovany706/url-shortener/internal/config"
	"github.com/rovany706/url-shortener/internal/handlers"
	"github.com/rovany706/url-shortener/internal/middleware"
	"github.com/rovany706/url-shortener/internal/repository"
	"go.uber.org/zap"
)

func MainRouter(app app.URLShortener, appConfig *config.AppConfig, repository repository.Repository, auth auth.JWTAuthentication, logger *zap.Logger) chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.ResponseLogger(logger))
	r.Use(middleware.RequestLogger(logger))
	r.Use(middleware.RequestGzipCompress())
	r.Use(middleware.ResponseGzipCompress())

	r.Route("/", func(r chi.Router) {
		r.Get("/{id}", handlers.RedirectHandler(app))
		r.Post("/", handlers.MakeShortURLHandler(app, auth, repository, appConfig))
		r.Get("/ping", handlers.PingHandler(repository, logger))
		r.Route("/api", func(r chi.Router) {
			r.Post("/shorten", handlers.MakeShortURLHandlerJSON(app, appConfig, auth, repository, logger))
			r.Post("/shorten/batch", handlers.MakeShortURLBatchHandler(app, appConfig, auth, repository, logger))
			r.Get("/user/urls", handlers.GetUserURLs(auth, repository, appConfig, logger))
		})
	})

	return r
}
