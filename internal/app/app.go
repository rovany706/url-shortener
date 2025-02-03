package app

import (
	"crypto/sha1"
	"fmt"
	"net/url"

	"github.com/rovany706/url-shortener/internal/repository"
)

const shortHashByteCount = 4

type URLShortener interface {
	GetFullURL(shortID string) (fullURL string, ok bool)
	GetShortID(fullURL string) (shortID string, err error)
}

type URLShortenerApp struct {
	repository repository.Repository
}

// Создает экземпляр URLShortenerApp с загруженными значениями
func NewURLShortenerApp(repository repository.Repository) *URLShortenerApp {
	app := URLShortenerApp{
		repository: repository,
	}

	return &app
}

// Метод GetFullURL возвращает полную ссылку по короткому id и флаг успеха операции.
func (app *URLShortenerApp) GetFullURL(shortID string) (fullURL string, ok bool) {
	return app.repository.GetFullURL(shortID)
}

// Метод GetShortID возвращает первые 4 байта sha1-хеша ссылки в виде строки.
func (app *URLShortenerApp) GetShortID(fullURL string) (shortID string, err error) {
	if _, err = url.ParseRequestURI(fullURL); err != nil {
		return "", err
	}

	hash := sha1.Sum([]byte(fullURL))
	shortHash := hash[:shortHashByteCount]
	shortID = fmt.Sprintf("%x", shortHash)

	if err := app.repository.SaveEntry(shortID, fullURL); err != nil {
		return "", err
	}

	return shortID, nil
}
