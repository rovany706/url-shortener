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

// URLShortener интерфейс сокращателя ссылок
type URLShortener interface {
	GetFullURL(ctx context.Context, shortID string) (shortenedURLInfo *repository.ShortenedURLInfo, ok bool)
	GetShortID(ctx context.Context, userID int, fullURL string) (shortID string, err error)
	GetShortIDBatch(ctx context.Context, userID int, fullURLs []string) (shortIDs []string, err error)
}

// URLShortenerApp реализует интерфейс URLShortener
type URLShortenerApp struct {
	repository repository.Repository
}

// NewURLShortenerApp создает экземпляр URLShortenerApp
func NewURLShortenerApp(repository repository.Repository) *URLShortenerApp {
	app := URLShortenerApp{
		repository: repository,
	}

	return &app
}

// GetFullURL возвращает полную ссылку по короткому id и флаг успеха операции.
func (app *URLShortenerApp) GetFullURL(ctx context.Context, shortID string) (shortenedURLInfo *repository.ShortenedURLInfo, ok bool) {
	return app.repository.GetFullURL(ctx, shortID)
}

// GetShortID возвращает первые 4 байта sha1-хеша ссылки в виде строки.
func (app *URLShortenerApp) GetShortID(ctx context.Context, userID int, fullURL string) (shortID string, err error) {
	if _, err = url.ParseRequestURI(fullURL); err != nil {
		return "", err
	}

	shortID = getShortSHA1Hash(fullURL, shortHashByteCount)

	if err := app.repository.SaveEntry(ctx, userID, shortID, fullURL); err != nil {
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

// GetShortIDBatch возвращает короткие ID слайса ссылок.
func (app *URLShortenerApp) GetShortIDBatch(ctx context.Context, userID int, fullURLs []string) (shortIDs []string, err error) {
	shortIDs = make([]string, len(fullURLs))
	for i, fullURL := range fullURLs {
		if _, err = url.ParseRequestURI(fullURL); err != nil {
			return nil, err
		}

		shortID := getShortSHA1Hash(fullURL, shortHashByteCount)
		shortIDs[i] = shortID
	}

	if err = app.saveBatch(ctx, userID, shortIDs, fullURLs); err != nil {
		return nil, err
	}

	return shortIDs, nil
}

func (app *URLShortenerApp) saveBatch(ctx context.Context, userID int, shortIDs []string, fullURLs []string) error {
	shortURLMap := make(map[string]string, len(shortIDs))

	for i := 0; i < len(shortIDs); i++ {
		shortURLMap[shortIDs[i]] = fullURLs[i]
	}

	return app.repository.SaveEntries(ctx, userID, shortURLMap)
}

func getShortSHA1Hash(value string, byteCount int) string {
	hash := sha1.Sum([]byte(value))
	shortHash := hash[:byteCount]

	return fmt.Sprintf("%x", shortHash)
}
