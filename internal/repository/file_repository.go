package repository

import (
	"context"
	"errors"
	"sync"

	"github.com/spf13/afero"

	"github.com/rovany706/url-shortener/internal/models"
	"github.com/rovany706/url-shortener/internal/storage"
)

// FileRepository репозиторий, использующий файл
type FileRepository struct {
	fs              afero.Fs
	storageFilepath string
	shortURLMap     *sync.Map
}

// NewFileRepository создает файл для хранения данных
func NewFileRepository(fs afero.Fs, storageFilepath string) (*FileRepository, error) {
	fileStorageReader, err := storage.NewFileStorageReader(fs, storageFilepath)

	if err != nil {
		return nil, err
	}

	defer func() {
		dErr := fileStorageReader.Close()
		err = errors.Join(err, dErr)
	}()

	storage, err := fileStorageReader.ReadAllEntries()

	if err != nil {
		return nil, err
	}

	repository := FileRepository{
		fs:              fs,
		storageFilepath: storageFilepath,
		shortURLMap:     initializeShortURLMap(storage),
	}

	return &repository, nil
}

func initializeShortURLMap(storage storage.Storage) *sync.Map {
	var shortURLMap sync.Map
	for _, v := range storage {
		shortURLMap.Store(v.ShortID, v.FullURL)
	}

	return &shortURLMap
}

// GetFullURL ищет в хранилище полную ссылку на ресурс по короткому ID
func (repository *FileRepository) GetFullURL(ctx context.Context, shortID string) (shortenedURLInfo *ShortenedURLInfo, ok bool) {
	v, ok := repository.shortURLMap.Load(shortID)

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
func (repository *FileRepository) SaveEntry(ctx context.Context, userID int, shortID string, fullURL string) error {
	_, exists := repository.shortURLMap.LoadOrStore(shortID, fullURL)

	if !exists {
		return repository.saveNewEntry(shortID, fullURL)
	}

	return nil
}

func (repository *FileRepository) saveNewEntry(shortID string, fullURL string) error {
	storageWriter, err := storage.NewFileStorageWriter(repository.fs, repository.storageFilepath)

	if err != nil {
		return err
	}

	defer func() {
		dErr := storageWriter.Close()
		err = errors.Join(err, dErr)
	}()

	entry := storage.StorageEntry{
		ShortID: shortID,
		FullURL: fullURL,
	}

	return storageWriter.WriteEntry(&entry)
}

// SaveEntries записывает набор сокращенных ссылок
func (repository *FileRepository) SaveEntries(ctx context.Context, userID int, shortIDMap URLMapping) error {
	storageWriter, err := storage.NewFileStorageWriter(repository.fs, repository.storageFilepath)

	if err != nil {
		return err
	}

	defer func() {
		dErr := storageWriter.Close()
		err = errors.Join(err, dErr)
	}()

	entries := make([]storage.StorageEntry, 0)
	for shortID, fullURL := range shortIDMap {
		_, exists := repository.shortURLMap.LoadOrStore(shortID, fullURL)

		if !exists {
			entry := storage.StorageEntry{
				ShortID: shortID,
				FullURL: fullURL,
			}

			entries = append(entries, entry)
		}
	}

	return storageWriter.WriteEntries(entries)
}

// Close завершает работу с хранилищем
func (repository *FileRepository) Close() error {
	return nil
}

// Ping не поддерживается FileRepository
func (repository *FileRepository) Ping(ctx context.Context) error {
	return ErrPingNotSupported
}

// GetShortID возвращает shortID сокращенной ссылки
func (repository *FileRepository) GetShortID(ctx context.Context, fullURL string) (shortID string, err error) {
	repository.shortURLMap.Range(func(id, url any) bool {
		if url.(string) == fullURL {
			shortID = id.(string)
			return false
		}
		return true
	})

	return
}

// GetUserEntries возвращает сокращенный пользователем ссылки по userID
func (repository *FileRepository) GetUserEntries(ctx context.Context, userID int) (shortIDMap URLMapping, err error) {
	return nil, ErrNotImplemented
}

// GetNewUserID возвращает ID нового пользователя
func (repository *FileRepository) GetNewUserID(ctx context.Context) (userID int, err error) {
	return -1, nil
}

// DeleteUserURLs удаляет набор сокращенных ссылок
func (repository *FileRepository) DeleteUserURLs(ctx context.Context, deleteRequests []models.UserDeleteRequest) error {
	return nil
}
