package app

import (
	"crypto/sha1"
	"fmt"
	"net/url"
	"sync"

	"github.com/rovany706/url-shortener/internal/storage"
)

const shortHashByteCount = 4

type URLShortener interface {
	GetFullURL(shortID string) (fullURL string, ok bool)
	GetShortID(fullURL string) (shortID string, err error)
}

type URLShortenerApp struct {
	storageFilepath string
	shortURLMap     sync.Map
}

// Создает экземпляр URLShortenerApp с загруженными значениями
func NewURLShortenerApp(storage storage.Storage) *URLShortenerApp {
	app := URLShortenerApp{}
	app.initializeShortURLMap(storage)

	return &app
}

func (app *URLShortenerApp) initializeShortURLMap(storage storage.Storage) {
	for _, v := range storage {
		app.shortURLMap.Store(v.ShortID, v.FullURL)
	}
}

func (app *URLShortenerApp) UseStorageFromFile(filename string) error {
	fileStorageReader, err := storage.NewFileStorageReader(filename)

	if err != nil {
		return err
	}

	defer fileStorageReader.Close()

	storage, err := fileStorageReader.ReadAllEntries()

	if err != nil {
		return err
	}

	app.initializeShortURLMap(storage)
	app.storageFilepath = filename

	return nil
}

// Метод GetFullURL возвращает полную ссылку по короткому id и флаг успеха операции.
func (app *URLShortenerApp) GetFullURL(shortID string) (fullURL string, ok bool) {
	v, ok := app.shortURLMap.Load(shortID)

	if ok {
		fullURL = v.(string)
	}

	return fullURL, ok
}

// Метод GetShortID возвращает первые 4 байта sha1-хеша ссылки в виде строки.
func (app *URLShortenerApp) GetShortID(fullURL string) (shortID string, err error) {
	if _, err = url.ParseRequestURI(fullURL); err != nil {
		return "", err
	}

	hash := sha1.Sum([]byte(fullURL))
	shortHash := hash[:shortHashByteCount]
	shortID = fmt.Sprintf("%x", shortHash)

	_, exists := app.shortURLMap.LoadOrStore(shortID, fullURL)

	if !exists {
		app.saveNewEntry(shortID, fullURL)
	}

	return shortID, err
}

func (app *URLShortenerApp) saveNewEntry(shortID string, fullURL string) error {
	if app.storageFilepath == "" {
		return nil
	}

	storageWriter, err := storage.NewFileStorageWriter(app.storageFilepath)

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
