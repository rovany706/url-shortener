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

func (r *MemoryRepository) Close() error {
	return nil
}

func (r *MemoryRepository) Ping(ctx context.Context) error {
	return ErrPingNotSupported
}

func (r *MemoryRepository) SaveEntries(ctx context.Context, shortIDMap map[string]string) error {
	for shortID, fullURL := range shortIDMap {
		err := r.SaveEntry(ctx, shortID, fullURL)

		if err != nil {
			return err
		}
	}

	return nil
}

func (r *MemoryRepository) GetShortID(ctx context.Context, fullURL string) (shortID string, err error) {
	r.shortURLMap.Range(func(id, url any) bool {
		if url.(string) == fullURL {
			shortID = id.(string)
			return false
		}
		return true
	})

	return
}
