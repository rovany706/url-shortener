package main

import (
	"os"

	"go.uber.org/zap"

	"github.com/rovany706/url-shortener/internal/config"
	"github.com/rovany706/url-shortener/internal/logger"
	"github.com/rovany706/url-shortener/internal/server"
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

	server, err := server.NewServer(appConfig, logger)
	if err != nil {
		logger.Fatal("error when creating server", zap.Error(err))
	}

	if err = run(server); err != nil {
		logger.Fatal("error when running server", zap.Error(err))
	}
}

func run(server *server.Server) error {
	defer server.StopServer()

	return server.RunServer()
}
