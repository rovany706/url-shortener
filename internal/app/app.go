package app

import (
	"crypto/sha1"
	"fmt"
	"net/url"
)

const shortHashByteCount = 4

var shortURLMap map[string]string

func init() {
	shortURLMap = make(map[string]string, 0)
}

// Метод GetFullURL возвращает полную ссылку по короткому id и флаг успеха оперции.
func GetFullURL(shortID string) (fullURL string, ok bool) {
	fullURL, ok = shortURLMap[shortID]
	return fullURL, ok
}

// Метод GetShortID возвращает первые 4 байта sha1-хеша ссылки в виде строки.
func GetShortID(fullURL string) (shortID string, err error) {
	if _, err = url.ParseRequestURI(fullURL); err != nil {
		return "", err
	}

	hash := sha1.Sum([]byte(fullURL))
	shortHash := hash[:shortHashByteCount]
	shortID = fmt.Sprintf("%x", shortHash)
	shortURLMap[shortID] = fullURL

	return shortID, err
}
