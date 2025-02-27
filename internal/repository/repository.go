package repository

import (
	"context"
	"errors"

	"github.com/rovany706/url-shortener/internal/config"
	"github.com/spf13/afero"
)

type ShortIDMap map[string]string

type Repository interface {
	GetFullURL(ctx context.Context, shortID string) (fullURL string, ok bool)
	SaveEntry(ctx context.Context, userID int, shortID string, fullURL string) error
	SaveEntries(ctx context.Context, userID int, shortIDMap ShortIDMap) error
	GetShortID(ctx context.Context, fullURL string) (shortID string, err error)
	GetUserEntries(ctx context.Context, userID int) (shortIDMap ShortIDMap, err error)
	GetNewUserID(ctx context.Context) (userID int, err error)
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
