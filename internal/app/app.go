package app

import (
	"context"
	"crypto/sha1"
	"errors"
	"fmt"
	"net/url"

	"github.com/rovany706/url-shortener/internal/repository"
)

const shortHashByteCount = 4

type URLShortener interface {
	GetFullURL(ctx context.Context, shortID string) (fullURL string, ok bool)
	GetShortID(ctx context.Context, fullURL string) (shortID string, err error)
	GetShortIDBatch(ctx context.Context, fullURLs []string) (shortIDs []string, err error)
}

type URLShortenerApp struct {
	repository repository.Repository
}

// Создает экземпляр URLShortenerApp
func NewURLShortenerApp(repository repository.Repository) *URLShortenerApp {
	app := URLShortenerApp{
		repository: repository,
	}

	return &app
}

// Метод GetFullURL возвращает полную ссылку по короткому id и флаг успеха операции.
func (app *URLShortenerApp) GetFullURL(ctx context.Context, shortID string) (fullURL string, ok bool) {
	return app.repository.GetFullURL(ctx, shortID)
}

// Метод GetShortID возвращает первые 4 байта sha1-хеша ссылки в виде строки.
func (app *URLShortenerApp) GetShortID(ctx context.Context, fullURL string) (shortID string, err error) {
	if _, err = url.ParseRequestURI(fullURL); err != nil {
		return "", err
	}

	shortID = getShortSHA1Hash(fullURL, shortHashByteCount)

	if err := app.repository.SaveEntry(ctx, shortID, fullURL); err != nil {
		switch {
		case errors.Is(err, repository.ErrConflict):
			shortID, err = app.repository.GetShortID(ctx, fullURL)

			if err != nil {
				return "", err
			}

			return shortID, repository.ErrConflict
		default:
			return "", err
		}
	}

	return shortID, nil
}

// Метод GetShortIDBatch возвращает короткие ID слайса ссылок.
func (app *URLShortenerApp) GetShortIDBatch(ctx context.Context, fullURLs []string) (shortIDs []string, err error) {
	shortIDs = make([]string, len(fullURLs))
	for i, fullURL := range fullURLs {
		if _, err = url.ParseRequestURI(fullURL); err != nil {
			return nil, err
		}

		shortID := getShortSHA1Hash(fullURL, shortHashByteCount)
		shortIDs[i] = shortID
	}

	if err = app.saveBatch(ctx, shortIDs, fullURLs); err != nil {
		return nil, err
	}

	return shortIDs, nil
}

func (app *URLShortenerApp) saveBatch(ctx context.Context, shortIDs []string, fullURLs []string) error {
	shortURLMap := make(map[string]string, len(shortIDs))

	for i := 0; i < len(shortIDs); i++ {
		shortURLMap[shortIDs[i]] = fullURLs[i]
	}

	return app.repository.SaveEntries(ctx, shortURLMap)
}

func getShortSHA1Hash(value string, byteCount int) string {
	hash := sha1.Sum([]byte(value))
	shortHash := hash[:byteCount]

	return fmt.Sprintf("%x", shortHash)
}
