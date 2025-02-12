package app

import (
	"errors"
	"strconv"
)

// Valid shortener mock
type MockURLShortener struct {
	shortURLMap map[string]string
	counter     int
}

func (shortener *MockURLShortener) GetFullURL(shortID string) (fullURL string, ok bool) {
	fullURL, ok = shortener.shortURLMap[shortID]
	return fullURL, ok
}

func (shortener *MockURLShortener) GetShortID(fullURL string) (shortID string, err error) {
	shortID = strconv.Itoa(shortener.counter)
	shortener.shortURLMap[shortID] = fullURL
	shortener.counter++

	return shortID, err
}

func NewMockURLShortener(shortURLMap map[string]string) *MockURLShortener {
	return &MockURLShortener{
		shortURLMap: shortURLMap,
	}
}

// Error shortener mock
type ErrMockURLShortener struct{}

func (shortener *ErrMockURLShortener) GetFullURL(shortID string) (fullURL string, ok bool) {
	return "", false
}

func (shortener *ErrMockURLShortener) GetShortID(fullURL string) (shortID string, err error) {
	return "", errors.New("test error")
}
