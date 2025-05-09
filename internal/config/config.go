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

// Ошибки
var (
	// ErrInvalidBaseURL ошибка валидации базового URL
	ErrInvalidBaseURL = errors.New("invalid base URL")
	// ErrInvalidAppRunAddress ошибка валидации адреса для запуска сервера
	ErrInvalidAppRunAddress = errors.New("invalid address and port to run server")
	// ErrInvalidLogLevel ошибка валидации уровня логгирования
	ErrInvalidLogLevel = errors.New("invalid log level")
)

const (
	defaultBaseURL         = "http://localhost:8080"
	defaultAppRunAddress   = ":8080"
	defaultLogLevel        = "info"
	defaultFileStoragePath = ""
	defaultDatabaseDSN     = ""
	defaultProfiling       = false
)

// StorageType тип хранилища данных сервиса
type StorageType int

// Перечисление типов хранилища
const (
	// None значение по умолчанию
	None StorageType = iota
	// File файловое хранилище
	File
	// Database хранение в базе данных
	Database
)

// AppConfig содержит конфигурацию сервиса
type AppConfig struct {
	// BaseURL базовый URL для сокращенных ссылок
	BaseURL string `env:"BASE_URL"`
	// AppRunAddress адрес для запуска сервера
	AppRunAddress string `env:"SERVER_ADDRESS"`
	// LogLevel уровень логгирования
	LogLevel string `env:"LOG_LEVEL"`
	// FileStoragePath путь файлового хранилища
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	// DatabaseDSN строка подключения к БД
	DatabaseDSN string `env:"DATABASE_DSN"`
	// EnableProfiling флаг включения режима профилирования
	EnableProfiling bool `env:"PPROF"`
	// StorageType тип хранилища
	StorageType StorageType
}

// Option функциональная опция
type Option func(*AppConfig)

// WithBaseURL задает базовый URL для сокращенных ссылок
func WithBaseURL(url string) Option {
	return func(c *AppConfig) {
		if url != "" {
			c.BaseURL = url
		}
	}
}

// WithAppRunAddress задает адрес для запуска сервера
func WithAppRunAddress(appRunAddress string) Option {
	return func(c *AppConfig) {
		if appRunAddress != "" {
			c.AppRunAddress = appRunAddress
		}
	}
}

// WithLogLevel задает уровень логгирования
func WithLogLevel(logLevel string) Option {
	return func(c *AppConfig) {
		if logLevel != "" {
			c.LogLevel = logLevel
		}
	}
}

// WithFileStoragePath задает путь файлового хранилища
func WithFileStoragePath(fileStoragePath string) Option {
	return func(c *AppConfig) {
		if fileStoragePath != "" {
			c.FileStoragePath = fileStoragePath
		}
	}
}

// WithDatabseDSN задает строку подключения к БД
func WithDatabseDSN(databaseDSN string) Option {
	return func(c *AppConfig) {
		if databaseDSN != "" {
			c.DatabaseDSN = databaseDSN
		}
	}
}

// WithStorageType задает тип хранилища
func WithStorageType(storageType StorageType) Option {
	return func(c *AppConfig) {
		c.StorageType = storageType
	}
}

// WithProfiling задает режим профилирования
func WithProfiling() Option {
	return func(c *AppConfig) {
		c.EnableProfiling = true
	}
}

// NewConfig создает экземпляр конфига.
// Поля настраиваются методами With...()
func NewConfig(opts ...Option) *AppConfig {
	cfg := &AppConfig{
		BaseURL:         defaultBaseURL,
		AppRunAddress:   defaultAppRunAddress,
		LogLevel:        defaultLogLevel,
		FileStoragePath: defaultFileStoragePath,
		DatabaseDSN:     defaultDatabaseDSN,
	}

	for _, opt := range opts {
		opt(cfg)
	}

	return cfg
}

// ParseArgs парсит параметры командной строки и возвращает указатель на объект AppConfig с заполненными значениями конфигурации
func ParseArgs(programName string, args []string) (appConfig *AppConfig, err error) {
	appConfig = new(AppConfig)
	flags := flag.NewFlagSet(programName, flag.ExitOnError)

	flags.StringVar(&appConfig.AppRunAddress, "a", defaultAppRunAddress, fmt.Sprintf("address and port to run server (default: %s)", defaultAppRunAddress))
	flags.StringVar(&appConfig.BaseURL, "b", defaultBaseURL, fmt.Sprintf("base URL for short links (default: %s)", defaultBaseURL))
	flags.StringVar(&appConfig.LogLevel, "l", defaultLogLevel, fmt.Sprintf("log level (default: %s)", defaultLogLevel))
	flags.StringVar(&appConfig.FileStoragePath, "f", defaultFileStoragePath, "file storage path")
	flags.StringVar(&appConfig.DatabaseDSN, "d", defaultDatabaseDSN, "database DSN")
	flags.BoolVar(&appConfig.EnableProfiling, "p", defaultProfiling, "enable pprof server at /debug")

	err = flags.Parse(args)

	if err != nil {
		return nil, err
	}

	err = env.Parse(appConfig)

	if err != nil {
		return nil, err
	}

	if err := validateParsedArgs(appConfig); err != nil {
		return nil, err
	}

	appConfig.StorageType = getStorageType(appConfig)

	log.Printf("Parsed app config: %+v\n", appConfig)

	return appConfig, nil
}

func validateParsedArgs(appConfig *AppConfig) error {
	if ok := isURL(appConfig.BaseURL); !ok {
		return ErrInvalidBaseURL
	}

	if _, err := net.ResolveTCPAddr("tcp", appConfig.AppRunAddress); err != nil {
		return ErrInvalidAppRunAddress
	}

	if _, err := zap.ParseAtomicLevel(appConfig.LogLevel); err != nil {
		return ErrInvalidLogLevel
	}

	return nil
}

func isURL(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}

func getStorageType(appConfig *AppConfig) StorageType {
	if appConfig.DatabaseDSN != "" {
		return Database
	} else if appConfig.FileStoragePath != "" {
		return File
	}

	return None
}
