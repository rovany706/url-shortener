package app

import (
	"crypto/sha1"
	"fmt"
	"net/url"
	"sync"
)

const shortHashByteCount = 4

type URLShortener interface {
	GetFullURL(shortID string) (fullURL string, ok bool)
	GetShortID(fullURL string) (shortID string, err error)
}

type URLShortenerApp struct {
	ShortURLMap sync.Map
}

// Метод GetFullURL возвращает полную ссылку по короткому id и флаг успеха операции.
func (app *URLShortenerApp) GetFullURL(shortID string) (fullURL string, ok bool) {
	v, ok := app.ShortURLMap.Load(shortID)

	if ok {
		fullURL = v.(string)
	}

	return fullURL, ok
}

// Метод GetShortID возвращает первые 4 байта sha1-хеша ссылки в виде строки.
func (app *URLShortenerApp) GetShortID(fullURL string) (shortID string, err error) {
	if _, err = url.ParseRequestURI(fullURL); err != nil {
		return "", err
	}

	hash := sha1.Sum([]byte(fullURL))
	shortHash := hash[:shortHashByteCount]
	shortID = fmt.Sprintf("%x", shortHash)
	app.ShortURLMap.Store(shortID, fullURL)

	return shortID, err
}
