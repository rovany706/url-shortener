package config

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"net/url"

	"github.com/caarlos0/env/v11"
	"go.uber.org/zap"
)

const ErrInvalidBaseURL = "invalid base URL"
const ErrInvalidAppRunAddress = "invalid address and port to run server"
const ErrInvalidLogLevel = "invalid log level"

type AppConfig struct {
	BaseURL       string `env:"BASE_URL"`
	AppRunAddress string `env:"SERVER_ADDRESS"`
	LogLevel      string `env:"LOG_LEVEL"`
}

const defaultBaseURL = "http://localhost:8080"
const defaultAppRunAddress = ":8080"
const defaultLogLevel = "info"

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

func WithLogLevel(logLevel string) Option {
	return func(c *AppConfig) {
		if logLevel != "" {
			c.LogLevel = logLevel
		}
	}
}

func NewConfig(opts ...Option) *AppConfig {
	cfg := &AppConfig{
		BaseURL:       defaultBaseURL,
		AppRunAddress: defaultAppRunAddress,
		LogLevel:      defaultLogLevel,
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
	flags.StringVar(&appConfig.LogLevel, "l", defaultLogLevel, fmt.Sprintf("log level (default: %s)", defaultLogLevel))

	err = flags.Parse(args)

	if err != nil {
		return nil, err
	}

	err = env.Parse(appConfig)

	if err != nil {
		return nil, err
	}

	log.Printf("Parsed app config: %+v\n", appConfig)

	if err := validateParsedArgs(appConfig); err != nil {
		return nil, err
	}

	return appConfig, nil
}

func validateParsedArgs(appConfig *AppConfig) error {
	if ok := isURL(appConfig.BaseURL); !ok {
		return errors.New(ErrInvalidBaseURL)
	}

	if _, err := net.ResolveTCPAddr("tcp", appConfig.AppRunAddress); err != nil {
		return errors.New(ErrInvalidAppRunAddress)
	}

	if _, err := zap.ParseAtomicLevel(appConfig.LogLevel); err != nil {
		return errors.New(ErrInvalidLogLevel)
	}

	return nil
}

func isURL(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}
