package config

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"net/url"

	"github.com/caarlos0/env/v11"
)

const ErrInvalidBaseURL = "invalid base URL"
const ErrInvalidAppRunAddress = "invalid address and port to run server"

type AppConfig struct {
	BaseURL       string `env:"BASE_URL"`
	AppRunAddress string `env:"SERVER_ADDRESS"`
}

const defaultBaseURL = "http://localhost:8080"
const defaultAppRunAddress = ":8080"

type Option func(*AppConfig)

func WithBaseURL(url string) Option {
	return func(c *AppConfig) {
		if url != "" {
			c.BaseURL = url
		}
	}
}

func WithAppRunAddress(appRunAddress string) Option {
	return func(c *AppConfig) {
		if appRunAddress != "" {
			c.AppRunAddress = appRunAddress
		}
	}
}

func NewConfig(opts ...Option) *AppConfig {
	cfg := &AppConfig{
		BaseURL:       defaultBaseURL,
		AppRunAddress: defaultAppRunAddress,
	}

	for _, opt := range opts {
		opt(cfg)
	}

	return cfg
}

// Метод ParseArgs парсит параметры командной строки и возвращает указатель на объект AppConfig с заполненными значениями конфигурации
func ParseArgs(programName string, args []string) (appConfig *AppConfig, err error) {
	appConfig = new(AppConfig)
	flags := flag.NewFlagSet(programName, flag.ExitOnError)

	flags.StringVar(&appConfig.AppRunAddress, "a", defaultAppRunAddress, fmt.Sprintf("address and port to run server (default: %s)", defaultAppRunAddress))
	flags.StringVar(&appConfig.BaseURL, "b", defaultBaseURL, fmt.Sprintf("base URL for short links (default: %s)", defaultBaseURL))

	err = flags.Parse(args)

	if err != nil {
		return nil, err
	}

	err = env.Parse(appConfig)

	if err != nil {
		return nil, err
	}

	log.Printf("Parsed app config: %+v\n", appConfig)

	return validateParsedArgs(appConfig)
}

func validateParsedArgs(appConfig *AppConfig) (*AppConfig, error) {
	if ok := isURL(appConfig.BaseURL); !ok {
		return nil, errors.New(ErrInvalidBaseURL)
	}

	if _, err := net.ResolveTCPAddr("tcp", appConfig.AppRunAddress); err != nil {
		return nil, errors.New(ErrInvalidAppRunAddress)
	}

	return appConfig, nil
}

func isURL(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}
