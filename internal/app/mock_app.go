package app

import (
	"context"
	"errors"
	"strconv"
)

// Valid shortener mock
type MockURLShortener struct {
	shortURLMap map[string]string
	counter     int
}

func (shortener *MockURLShortener) GetFullURL(ctx context.Context, shortID string) (fullURL string, ok bool) {
	fullURL, ok = shortener.shortURLMap[shortID]
	return fullURL, ok
}

func (shortener *MockURLShortener) GetShortID(ctx context.Context, userID int, fullURL string) (shortID string, err error) {
	shortID = strconv.Itoa(shortener.counter)
	shortener.shortURLMap[shortID] = fullURL
	shortener.counter++

	return shortID, nil
}

func (shortener *MockURLShortener) GetShortIDBatch(ctx context.Context, userID int, fullURLs []string) (shortIDs []string, err error) {
	shortIDs = make([]string, 0)
	for _, fullURL := range fullURLs {
		shortID, _ := shortener.GetShortID(ctx, userID, fullURL)
		shortIDs = append(shortIDs, shortID)
	}

	return shortIDs, nil
}

func NewMockURLShortener(shortURLMap map[string]string) *MockURLShortener {
	return &MockURLShortener{
		shortURLMap: shortURLMap,
	}
}

// Error shortener mock
type ErrMockURLShortener struct{}

func (shortener *ErrMockURLShortener) GetFullURL(ctx context.Context, shortID string) (fullURL string, ok bool) {
	return "", false
}

func (shortener *ErrMockURLShortener) GetShortID(ctx context.Context, userID int, fullURL string) (shortID string, err error) {
	return "", errors.New("test error")
}

func (shortener *ErrMockURLShortener) GetShortIDBatch(ctx context.Context, userID int, fullURLs []string) (shortIDs []string, err error) {
	return nil, errors.New("test error")
}
