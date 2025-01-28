package storage

import (
	"encoding/json"
	"os"
)

type FileStorageWriter struct {
	file    *os.File
	encoder *json.Encoder
}

func NewFileStorageWriter(filename string) (*FileStorageWriter, error) {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	return &FileStorageWriter{
		file:    file,
		encoder: json.NewEncoder(file),
	}, nil
}

func (w *FileStorageWriter) WriteEntry(entry *StorageEntry) error {
	return w.encoder.Encode(&entry)
}

func (w *FileStorageWriter) Close() error {
	return w.file.Close()
}

type FileStorageReader struct {
	file    *os.File
	decoder *json.Decoder
}

func NewFileStorageReader(filename string) (*FileStorageReader, error) {
	file, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	return &FileStorageReader{
		file:    file,
		decoder: json.NewDecoder(file),
	}, nil
}

func (r *FileStorageReader) ReadEntry() (*StorageEntry, error) {
	entry := &StorageEntry{}
	if err := r.decoder.Decode(&entry); err != nil {
		return nil, err
	}

	return entry, nil
}

func (r *FileStorageReader) ReadAllEntries() (Storage, error) {
	entries := make([]StorageEntry, 0)
	for r.decoder.More() {
		entry, err := r.ReadEntry()

		if err != nil {
			return nil, err
		}
		entries = append(entries, *entry)
	}

	return entries, nil
}

func (r *FileStorageReader) Close() error {
	return r.file.Close()
}
