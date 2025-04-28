package storage

import (
	"encoding/json"
	"os"

	"github.com/spf13/afero"
)

// FileStorageReader объект для чтения записей из файла
type FileStorageReader struct {
	file    afero.File
	decoder *json.Decoder
}

// NewFileStorageReader создает FileStorageReader
func NewFileStorageReader(fs afero.Fs, filename string) (*FileStorageReader, error) {
	file, err := fs.OpenFile(filename, os.O_RDONLY|os.O_CREATE, os.ModeType)
	if err != nil {
		return nil, err
	}

	return &FileStorageReader{
		file:    file,
		decoder: json.NewDecoder(file),
	}, nil
}

// ReadEntry возвращает запись
func (r *FileStorageReader) ReadEntry() (*StorageEntry, error) {
	entry := &StorageEntry{}
	if err := r.decoder.Decode(&entry); err != nil {
		return nil, err
	}

	return entry, nil
}

// ReadAllEntries возвращает все записи в файле
func (r *FileStorageReader) ReadAllEntries() (Storage, error) {
	entries := make(Storage, 0)
	for r.decoder.More() {
		entry, err := r.ReadEntry()

		if err != nil {
			return nil, err
		}
		entries = append(entries, *entry)
	}

	return entries, nil
}

// Close завершает работу с файлом
func (r *FileStorageReader) Close() error {
	return r.file.Close()
}
