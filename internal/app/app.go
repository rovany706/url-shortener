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
func GetFullURL(shortId string) (fullURL string, ok bool) {
	fullURL, ok = shortURLMap[shortId]
	return fullURL, ok
}

// Метод GetShortId возвращает первые 4 байта sha1-хеша ссылки в виде строки.
func GetShortId(fullURL string) (shortId string, err error) {
	if _, err = url.ParseRequestURI(fullURL); err != nil {
		return "", err
	}

	hash := sha1.Sum([]byte(fullURL))
	shortHash := hash[:shortHashByteCount]
	shortId = fmt.Sprintf("%x", shortHash)
	shortURLMap[shortId] = fullURL

	return shortId, err
}
