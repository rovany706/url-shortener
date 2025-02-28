package main

import (
	"context"
	"os"

	"github.com/rovany706/url-shortener/internal/app"
	"github.com/rovany706/url-shortener/internal/auth"
	"github.com/rovany706/url-shortener/internal/config"
	"github.com/rovany706/url-shortener/internal/logger"
	"github.com/rovany706/url-shortener/internal/repository"
	"github.com/rovany706/url-shortener/internal/server"
	"github.com/rovany706/url-shortener/internal/service"
	"go.uber.org/zap"
)

var appConfig *config.AppConfig

func main() {
	var err error
	appConfig, err = config.ParseArgs(os.Args[0], os.Args[1:])

	if err != nil {
		panic(err)
	}

	logger, err := logger.NewLogger(appConfig.LogLevel)

	if err != nil {
		panic(err)
	}

	if err = run(appConfig, logger); err != nil {
		logger.Fatal("error when running server", zap.Error(err))
	}
}

func run(appConfig *config.AppConfig, logger *zap.Logger) error {
	repository, err := repository.NewAppRepository(context.Background(), appConfig)

	if err != nil {
		logger.Fatal("failed to create repository", zap.Error(err))
	}

	auth, err := auth.NewAppJWTAuthentication(nil)

	if err != nil {
		logger.Fatal("failed to create authentication", zap.Error(err))
	}

	app := app.NewURLShortenerApp(repository)

	defer repository.Close()

	deleteService := service.NewDeleteService(repository)
	deleteService.StartWorker()
	defer deleteService.StopWorker()

	return server.RunServer(app, appConfig, auth, repository, deleteService, logger)
}
