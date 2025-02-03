package storage

import (
	"bufio"
	"encoding/json"
	"os"

	"github.com/spf13/afero"
)

type FileStorageWriter struct {
	file    afero.File
	encoder *json.Encoder
	buffer  *bufio.Writer
}

func NewFileStorageWriter(fs afero.Fs, filename string) (*FileStorageWriter, error) {
	file, err := fs.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	bufWriter := bufio.NewWriter(file)

	return &FileStorageWriter{
		file:    file,
		encoder: json.NewEncoder(bufWriter),
		buffer:  bufWriter,
	}, nil
}

func (w *FileStorageWriter) WriteEntry(entry *StorageEntry) error {
	if err := w.encoder.Encode(&entry); err != nil {
		return err
	}

	return w.buffer.Flush()
}

func (w *FileStorageWriter) Close() error {
	return w.file.Close()
}
