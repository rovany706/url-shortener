package repository

import (
	"context"
	"errors"

	"github.com/rovany706/url-shortener/internal/config"
	"github.com/rovany706/url-shortener/internal/models"
	"github.com/spf13/afero"
)

type URLMapping map[string]string

type ShortenedURLInfo struct {
	UserID    int
	FullURL   string
	ShortID   string
	IsDeleted bool
}

type Repository interface {
	GetFullURL(ctx context.Context, shortID string) (shortenedURLInfo *ShortenedURLInfo, ok bool)
	SaveEntry(ctx context.Context, userID int, shortID string, fullURL string) error
	SaveEntries(ctx context.Context, userID int, shortIDMap URLMapping) error
	GetShortID(ctx context.Context, fullURL string) (shortID string, err error)
	GetUserEntries(ctx context.Context, userID int) (shortIDMap URLMapping, err error)
	GetNewUserID(ctx context.Context) (userID int, err error)
	DeleteUserURLs(ctx context.Context, deleteRequests []models.UserDeleteRequest) error
	Ping(ctx context.Context) error
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
