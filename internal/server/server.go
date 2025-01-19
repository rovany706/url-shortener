package server

import (
	"net/http"

	"github.com/rovany706/url-shortener/internal/app"
	"github.com/rovany706/url-shortener/internal/config"
	"github.com/rovany706/url-shortener/internal/router"
)

func RunServer(app app.URLShortener, appConfig *config.AppConfig) error {
	r := router.MainRouter(app, appConfig)

	return http.ListenAndServe(appConfig.AppRunAddress, r)
}
