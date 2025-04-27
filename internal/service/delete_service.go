package service

import (
	"context"
	"sync"
	"time"

	"github.com/rovany706/url-shortener/internal/models"
	"github.com/rovany706/url-shortener/internal/repository"
)

const (
	deleteFlushTimePeriod = time.Second * 10
)

type DeleteService interface {
	Put(deleteChan chan models.UserDeleteRequest)
	StartWorker(context.Context)
}

type DeleteServiceImpl struct {
	flushTicker  *time.Ticker
	deleteBuffer *DeleteRequestBuffer
	repo         repository.Repository
}

func NewDeleteService(repo repository.Repository) *DeleteServiceImpl {
	return &DeleteServiceImpl{
		deleteBuffer: NewDeleteBuffer(),
		repo:         repo,
	}
}

func (ds *DeleteServiceImpl) Put(deleteChan chan models.UserDeleteRequest) {
	ds.deleteBuffer.Add(deleteChan)
}

func (ds *DeleteServiceImpl) StartWorker(ctx context.Context) {
	ds.flushTicker = time.NewTicker(deleteFlushTimePeriod)

	go func() {
		for {
			select {
			case <-ds.flushTicker.C:
				deleteChs := ds.deleteBuffer.Flush()
				deleteFanInCh := fanIn(deleteChs...)

				deleteRequests := make([]models.UserDeleteRequest, 0)
				for request := range deleteFanInCh {
					deleteRequests = append(deleteRequests, request)
				}

				_ = ds.repo.DeleteUserURLs(context.Background(), deleteRequests)
			case <-ctx.Done():
				ds.flushTicker.Stop()
				return
			}
		}
	}()
}

func fanIn(deleteChs ...chan models.UserDeleteRequest) chan models.UserDeleteRequest {
	resultCh := make(chan models.UserDeleteRequest)
	var wg sync.WaitGroup

	for _, ch := range deleteChs {
		chClosure := ch
		wg.Add(1)

		go func() {
			defer wg.Done()

			for request := range chClosure {
				resultCh <- request
			}
		}()
	}

	go func() {
		wg.Wait()
		close(resultCh)
	}()

	return resultCh
}
