package main

import (
	"os"

	"github.com/rovany706/url-shortener/internal/app"
	"github.com/rovany706/url-shortener/internal/config"
	"github.com/rovany706/url-shortener/internal/server"
)

var appConfig *config.AppConfig

func main() {
	var err error
	appConfig, err = config.ParseArgs(os.Args[0], os.Args[1:])

	if err != nil {
		panic(err)
	}

	if err := run(appConfig); err != nil {
		panic(err)
	}
}

func run(appConfig *config.AppConfig) error {
	app := app.URLShortenerApp{}

	return server.RunServer(&app, appConfig)
}
