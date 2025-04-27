package repository

import (
	"context"
	"errors"

	"github.com/rovany706/url-shortener/internal/config"
	"github.com/rovany706/url-shortener/internal/models"
	"github.com/spf13/afero"
)

// URLMapping словарь идентификатора и сокращенной ссылки
type URLMapping map[string]string

// ShortenedURLInfo информация о сокращенной ссылке
type ShortenedURLInfo struct {
	// UserID идентификатор пользователя
	UserID int
	// FullURL сокращенная ссылка
	FullURL string
	// ShortID идентификатор сокращенной ссылки
	ShortID string
	// IsDeleted флаг удаленной ссылки
	IsDeleted bool
}

// Repository интерфейс работы с данными сервиса
type Repository interface {
	// GetFullURL возвращает сокращенную ссылку по shortID
	GetFullURL(ctx context.Context, shortID string) (shortenedURLInfo *ShortenedURLInfo, ok bool)
	// SaveEntry записывает сокращаемую ссылку по shortID
	SaveEntry(ctx context.Context, userID int, shortID string, fullURL string) error
	// SaveEntries записывает набор сокращенных ссылок
	SaveEntries(ctx context.Context, userID int, shortIDMap URLMapping) error
	// GetShortID возвращает shortID сокращенной ссылки
	GetShortID(ctx context.Context, fullURL string) (shortID string, err error)
	// GetUserEntries возвращает сокращенный пользователем ссылки по userID
	GetUserEntries(ctx context.Context, userID int) (shortIDMap URLMapping, err error)
	// GetNewUserID возвращает ID нового пользователя
	GetNewUserID(ctx context.Context) (userID int, err error)
	// DeleteUserURLs удаляет набор сокращенных ссылок
	DeleteUserURLs(ctx context.Context, deleteRequests []models.UserDeleteRequest) error
	// Ping проверяет подключение к источнику данных
	Ping(ctx context.Context) error
	// Close завершает работу с источником данных
	Close() error
}

var (
	ErrUnknownStorageType = errors.New("unknown storage type")
	ErrPingNotSupported   = errors.New("ping is not supported for this storage type")
	ErrConflict           = errors.New("entry conflict")
	ErrNotImplemented     = errors.New("method is not implemented")
)

func NewAppRepository(ctx context.Context, appConfig *config.AppConfig) (Repository, error) {
	switch appConfig.StorageType {
	case config.Database:
		return NewDatabaseRepository(ctx, appConfig.DatabaseDSN)
	case config.File:
		return NewFileRepository(afero.NewOsFs(), appConfig.FileStoragePath)
	case config.None:
		return NewMemoryRepository(), nil
	default:
		return nil, ErrUnknownStorageType
	}
}
