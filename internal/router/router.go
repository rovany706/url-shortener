package router

import (
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/rovany706/url-shortener/internal/handlers"
	"github.com/rovany706/url-shortener/internal/middleware"
	"github.com/rovany706/url-shortener/internal/repository"
)

// GetRouter возвращает роутер сервиса
func GetRouter(
	shortenHandlers handlers.ShortenURLHandlers,
	userHandlers handlers.UserHandlers,
	redirectHandlers handlers.RedirectHandlers,
	repository repository.Repository,
	logger *zap.Logger,
) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.ResponseLogger(logger))
	r.Use(middleware.RequestLogger(logger))
	r.Use(middleware.RequestGzipCompress(logger))
	r.Use(middleware.ResponseGzipCompress(logger))

	r.Get("/ping", handlers.PingHandler(repository, logger))

	registerShortenHandlers(r, shortenHandlers)
	registerUserHandlers(r, userHandlers)
	registerRedirectHandlers(r, redirectHandlers)

	return r
}

func registerRedirectHandlers(router chi.Router, redirectHandlers handlers.RedirectHandlers) {
	router.Get("/{id}", redirectHandlers.RedirectHandler())
}

func registerShortenHandlers(router chi.Router, shortenHandlers handlers.ShortenURLHandlers) {
	router.Post("/", shortenHandlers.MakeShortURLHandler())
	router.Post("/api/shorten", shortenHandlers.MakeShortURLHandlerJSON())
	router.Post("/api/shorten/batch", shortenHandlers.MakeShortURLBatchHandler())
}

func registerUserHandlers(router chi.Router, userHandlers handlers.UserHandlers) {
	router.Get("/api/user/urls", userHandlers.GetUserURLsHandler())
	router.Delete("/api/user/urls", userHandlers.DeleteUserURLsHandler())
}
