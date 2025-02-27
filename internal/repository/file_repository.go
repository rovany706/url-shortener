package repository

import (
	"context"
	"sync"

	"github.com/rovany706/url-shortener/internal/storage"
	"github.com/spf13/afero"
)

type FileRepository struct {
	fs              afero.Fs
	storageFilepath string
	shortURLMap     *sync.Map
}

func NewFileRepository(fs afero.Fs, storageFilepath string) (*FileRepository, error) {
	fileStorageReader, err := storage.NewFileStorageReader(fs, storageFilepath)

	if err != nil {
		return nil, err
	}

	defer fileStorageReader.Close()

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

// Метод GetFullURL ищет в хранилище полную ссылку на ресурс по короткому ID
func (repository *FileRepository) GetFullURL(ctx context.Context, shortID string) (fullURL string, ok bool) {
	v, ok := repository.shortURLMap.Load(shortID)

	if ok {
		fullURL = v.(string)
	}

	return fullURL, ok
}

// Метод SaveEntry сохраняет в хранилище информацию о сокращенной ссылке
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

	defer storageWriter.Close()

	entry := storage.StorageEntry{
		ShortID: shortID,
		FullURL: fullURL,
	}

	return storageWriter.WriteEntry(&entry)
}

func (repository *FileRepository) SaveEntries(ctx context.Context, userID int, shortIDMap ShortIDMap) error {
	storageWriter, err := storage.NewFileStorageWriter(repository.fs, repository.storageFilepath)

	if err != nil {
		return err
	}

	defer storageWriter.Close()

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

func (repository *FileRepository) Close() error {
	return nil
}

func (repository *FileRepository) Ping(ctx context.Context) error {
	return ErrPingNotSupported
}

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

func (repository *FileRepository) GetUserEntries(ctx context.Context, userID int) (shortIDMap ShortIDMap, err error) {
	return nil, ErrNotImplemented
}

func (repository *FileRepository) GetNewUserID(ctx context.Context) (userID int, err error) {
	return -1, ErrNotImplemented
}
