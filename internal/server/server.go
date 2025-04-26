package server

import (
	"context"
	"net/http"

	"github.com/rovany706/url-shortener/internal/app"
	"github.com/rovany706/url-shortener/internal/auth"
	"github.com/rovany706/url-shortener/internal/config"
	"github.com/rovany706/url-shortener/internal/handlers"
	"github.com/rovany706/url-shortener/internal/repository"
	"github.com/rovany706/url-shortener/internal/router"
	"github.com/rovany706/url-shortener/internal/service"
	"go.uber.org/zap"
)

type Server struct {
	appConfig     *config.AppConfig
	app           app.URLShortener
	repository    repository.Repository
	deleteService service.DeleteService
	tokenManager  auth.TokenManager
	logger        *zap.Logger
}

func NewServer(appConfig *config.AppConfig, logger *zap.Logger) (*Server, error) {
	repository, err := repository.NewAppRepository(context.Background(), appConfig)
	if err != nil {
		return nil, err
	}

	tokenManager, err := auth.NewJWTTokenManager(nil)

	if err != nil {
		return nil, err
	}

	app := app.NewURLShortenerApp(repository)

	deleteService := service.NewDeleteService(repository)

	return &Server{
		appConfig:     appConfig,
		app:           app,
		repository:    repository,
		deleteService: deleteService,
		tokenManager:  tokenManager,
		logger:        logger,
	}, nil
}

func (server *Server) RunServer() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	server.deleteService.StartWorker(ctx)

	userHandlers := handlers.NewUserHandlers(server.deleteService, server.tokenManager, server.repository, server.appConfig, server.logger)
	redirectHandlers := handlers.NewRedirectHandlers(server.app)
	shortenHandlers := handlers.NewShortenURLHandlers(server.app, server.tokenManager, server.repository, server.appConfig, server.logger)
	r := router.GetRouter(shortenHandlers, userHandlers, redirectHandlers, server.repository, server.logger)

	return http.ListenAndServe(server.appConfig.AppRunAddress, r)
}

func (server *Server) StopServer() {
	server.repository.Close()
}
