package server

import (
	"net/http"

	"github.com/rovany706/url-shortener/internal/app"
	"github.com/rovany706/url-shortener/internal/config"
	"github.com/rovany706/url-shortener/internal/router"
	"go.uber.org/zap"
)

func RunServer(app app.URLShortener, appConfig *config.AppConfig, logger *zap.Logger) error {
	r := router.MainRouter(app, appConfig, logger)

	return http.ListenAndServe(appConfig.AppRunAddress, r)
}
