package repository

import (
	"context"
	"sync"
)

type MemoryRepository struct {
	shortURLMap sync.Map
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{}
}

func (r *MemoryRepository) GetFullURL(ctx context.Context, shortID string) (fullURL string, ok bool) {
	v, ok := r.shortURLMap.Load(shortID)

	if ok {
		fullURL = v.(string)
	}

	return fullURL, ok
}

func (r *MemoryRepository) SaveEntry(ctx context.Context, shortID string, fullURL string) error {
	r.shortURLMap.LoadOrStore(shortID, fullURL)

	return nil
}

func (r *MemoryRepository) Close() {

}

func (r *MemoryRepository) Ping(ctx context.Context) error {
	return ErrPingNotSupported
}
