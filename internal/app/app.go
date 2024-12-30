package app

import (
	"crypto/sha1"
	"fmt"
	"net/url"
)

const shortHashByteCount = 4

type URLShortenerApp struct {
	ShortURLMap map[string]string
}

// Метод GetFullURL возвращает полную ссылку по короткому id и флаг успеха операции.
func (app *URLShortenerApp) GetFullURL(shortID string) (fullURL string, ok bool) {
	fullURL, ok = app.ShortURLMap[shortID]
	return fullURL, ok
}

// Метод GetShortID возвращает первые 4 байта sha1-хеша ссылки в виде строки.
func (app *URLShortenerApp) GetShortID(fullURL string) (shortID string, err error) {
	if _, err = url.ParseRequestURI(fullURL); err != nil {
		return "", err
	}

	if app.ShortURLMap == nil {
		app.ShortURLMap = map[string]string{}
	}

	hash := sha1.Sum([]byte(fullURL))
	shortHash := hash[:shortHashByteCount]
	shortID = fmt.Sprintf("%x", shortHash)
	app.ShortURLMap[shortID] = fullURL

	return shortID, err
}
