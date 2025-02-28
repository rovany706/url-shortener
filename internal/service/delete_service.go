package service

import (
	"context"
	"sync"
	"time"

	"github.com/rovany706/url-shortener/internal/models"
	"github.com/rovany706/url-shortener/internal/repository"
)

const (
	DeleteFlushTimePeriod = time.Second * 10
	FanInChanCapacity     = 10
	Timeout               = time.Second * 10
)

type DeleteService struct {
	flushTicker  *time.Ticker
	deleteBuffer *DeleteRequestBuffer
	repo         repository.Repository
}

func NewDeleteService(repo repository.Repository) *DeleteService {
	return &DeleteService{
		deleteBuffer: NewDeleteBuffer(),
		repo:         repo,
	}
}

func (ds *DeleteService) Put(deleteChan chan models.UserDeleteRequest) {
	ds.deleteBuffer.Add(deleteChan)
}

func (ds *DeleteService) StartWorker() {
	ds.flushTicker = time.NewTicker(DeleteFlushTimePeriod)

	go func() {
		for range ds.flushTicker.C {
			deleteChs := ds.deleteBuffer.Flush()
			deleteFanInCh := fanIn(deleteChs...)

			deleteRequests := make([]models.UserDeleteRequest, 0)
			for request := range deleteFanInCh {
				deleteRequests = append(deleteRequests, request)
			}

			err := ds.repo.DeleteUserURLs(context.Background(), deleteRequests)
			if err != nil {
				// TODO
			}
		}
	}()
}

func (ds *DeleteService) StopWorker() {
	ds.flushTicker.Stop()
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
