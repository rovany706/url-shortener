package repository

import (
	"context"
	"errors"

	"github.com/rovany706/url-shortener/internal/config"
	"github.com/spf13/afero"
)

type Repository interface {
	GetFullURL(ctx context.Context, shortID string) (fullURL string, ok bool)
	SaveEntry(ctx context.Context, shortID string, fullURL string) error
	SaveEntries(ctx context.Context, shortIDMap map[string]string) error
	Ping(ctx context.Context) error
	Close() error
}

var (
	ErrUnknownStorageType = errors.New("unknown storage type")
	ErrPingNotSupported   = errors.New("ping is not supported for this storage type")
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
