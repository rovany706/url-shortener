package server

import (
	"net/http"

	"github.com/rovany706/url-shortener/internal/app"
	"github.com/rovany706/url-shortener/internal/auth"
	"github.com/rovany706/url-shortener/internal/config"
	"github.com/rovany706/url-shortener/internal/repository"
	"github.com/rovany706/url-shortener/internal/router"
	"github.com/rovany706/url-shortener/internal/service"
	"go.uber.org/zap"
)

func RunServer(app app.URLShortener, appConfig *config.AppConfig, auth auth.JWTAuthentication, repository repository.Repository, deleteService *service.DeleteService, logger *zap.Logger) error {
	r := router.MainRouter(app, appConfig, repository, auth, deleteService, logger)

	return http.ListenAndServe(appConfig.AppRunAddress, r)
}
