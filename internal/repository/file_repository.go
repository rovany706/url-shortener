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
func (repository *FileRepository) SaveEntry(ctx context.Context, shortID string, fullURL string) error {
	_, exists := repository.shortURLMap.LoadOrStore(shortID, fullURL)

	if !exists {
		return repository.saveNewEntry(shortID, fullURL)
	}

	return nil
}

func (repository *FileRepository) saveNewEntry(shortID string, fullURL string) error {
	if repository.storageFilepath == "" {
		return nil
	}

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

func (repository *FileRepository) Close() {

}

func (repository *FileRepository) Ping(ctx context.Context) error {
	return ErrPingNotSupported
}
