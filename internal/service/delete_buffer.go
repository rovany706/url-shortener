package service

import (
	"sync"

	"github.com/rovany706/url-shortener/internal/models"
)

// DeleteRequestBuffer буфер для хранения запросов на удаление укороченных ссылок
type DeleteRequestBuffer struct {
	buffer []chan models.UserDeleteRequest
	mutex  sync.RWMutex
}

// NewDeleteBuffer создает экземпляр DeleteRequestBuffer
func NewDeleteBuffer() *DeleteRequestBuffer {
	return &DeleteRequestBuffer{
		buffer: make([]chan models.UserDeleteRequest, 0),
	}
}

// Add добавляет канал с запросом на удаление в буфер
func (db *DeleteRequestBuffer) Add(request chan models.UserDeleteRequest) {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	db.buffer = append(db.buffer, request)
}

// Flush возвращает накопленные запросы и очищает буфер
func (db *DeleteRequestBuffer) Flush() []chan models.UserDeleteRequest {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	requests := append(make([]chan models.UserDeleteRequest, 0, len(db.buffer)), db.buffer...)
	db.buffer = nil

	return requests
}
