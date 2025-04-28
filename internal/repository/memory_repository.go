package repository

import (
	"context"
	"sync"

	"github.com/rovany706/url-shortener/internal/models"
)

// MemoryRepository репозиторий, хранящий информацию в памяти
type MemoryRepository struct {
	shortURLMap sync.Map
}

// NewMemoryRepository инициализирует работу с хранилищем в памяти
func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{}
}

// GetFullURL ищет в хранилище полную ссылку на ресурс по короткому ID
func (r *MemoryRepository) GetFullURL(ctx context.Context, shortID string) (shortenedURLInfo *ShortenedURLInfo, ok bool) {
	v, ok := r.shortURLMap.Load(shortID)

	if ok {
		fullURL := v.(string)
		shortenedURLInfo = &ShortenedURLInfo{
			UserID:    1,
			FullURL:   fullURL,
			ShortID:   shortID,
			IsDeleted: false,
		}
	}

	return shortenedURLInfo, ok
}

// SaveEntry сохраняет в хранилище информацию о сокращенной ссылке
func (r *MemoryRepository) SaveEntry(ctx context.Context, userID int, shortID string, fullURL string) error {
	r.shortURLMap.LoadOrStore(shortID, fullURL)

	return nil
}

// Close завершает работу с хранилищем
func (r *MemoryRepository) Close() error {
	return nil
}

// Ping не поддерживается MemoryRepository
func (r *MemoryRepository) Ping(ctx context.Context) error {
	return ErrPingNotSupported
}

// SaveEntries записывает набор сокращенных ссылок
func (r *MemoryRepository) SaveEntries(ctx context.Context, userID int, shortIDMap URLMapping) error {
	for shortID, fullURL := range shortIDMap {
		err := r.SaveEntry(ctx, userID, shortID, fullURL)

		if err != nil {
			return err
		}
	}

	return nil
}

// GetShortID возвращает shortID сокращенной ссылки
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

// GetUserEntries возвращает сокращенный пользователем ссылки по userID
func (r *MemoryRepository) GetUserEntries(ctx context.Context, userID int) (shortIDMap URLMapping, err error) {
	return nil, ErrNotImplemented
}

// GetNewUserID возвращает ID нового пользователя
func (r *MemoryRepository) GetNewUserID(ctx context.Context) (userID int, err error) {
	return -1, nil
}

// DeleteUserURLs удаляет набор сокращенных ссылок
func (r *MemoryRepository) DeleteUserURLs(ctx context.Context, deleteRequests []models.UserDeleteRequest) error {
	return nil
}
