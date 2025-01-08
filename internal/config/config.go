package config

import (
	"flag"
	"log"
)

type AppConfig struct {
	BaseURL       string
	AppRunAddress string
}

const defaultBaseURL = "http://localhost:8080"
const defaultAppRunAddress = ":8080"

// Метод ParseArgs парсит параметры командной строки и возвращает указатель на объект AppConfig с заполненными значениями конфигурации
func ParseArgs(programName string, args []string) (appConfig *AppConfig, err error) {
	appConfig = new(AppConfig)
	flags := flag.NewFlagSet(programName, flag.ExitOnError)

	flags.StringVar(&appConfig.AppRunAddress, "a", defaultAppRunAddress, "address and port to run server")
	flags.StringVar(&appConfig.BaseURL, "b", defaultBaseURL, "base URL for short links")

	err = flags.Parse(args)

	if err != nil {
		return nil, err
	}

	log.Printf("Parsed app config: %+v\n", appConfig)

	return appConfig, nil
}
