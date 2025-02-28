package service

import (
	"sync"

	"github.com/rovany706/url-shortener/internal/models"
)

type DeleteRequestBuffer struct {
	buffer []chan models.UserDeleteRequest
	mutex  sync.RWMutex
}

func NewDeleteBuffer() *DeleteRequestBuffer {
	return &DeleteRequestBuffer{
		buffer: make([]chan models.UserDeleteRequest, 0),
	}
}

func (db *DeleteRequestBuffer) Add(request chan models.UserDeleteRequest) {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	db.buffer = append(db.buffer, request)
}

func (db *DeleteRequestBuffer) Flush() []chan models.UserDeleteRequest {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	requests := append(make([]chan models.UserDeleteRequest, 0, len(db.buffer)), db.buffer...)
	db.buffer = nil

	return requests
}

func (db *DeleteRequestBuffer) Len() int {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	return len(db.buffer)
}
