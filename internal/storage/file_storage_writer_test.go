package storage

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

func TestNewFileStorageWriter(t *testing.T) {
	fs := afero.NewMemMapFs()
	filePath := "home/test/storage.json"
	_, err := NewFileStorageWriter(fs, filePath)
	require.NoError(t, err)

	stat, err := fs.Stat(filePath)
	require.NoError(t, err)
	require.NotEmpty(t, stat)
}

func TestWriteEntry(t *testing.T) {
	fs := afero.NewMemMapFs()
	filePath := "home/test/storage.json"
	writer, err := NewFileStorageWriter(fs, filePath)
	require.NoError(t, err)

	entry := StorageEntry{
		ShortID: "1",
		FullURL: "http://example.com",
	}
	err = writer.WriteEntry(&entry)
	require.NoError(t, err)

	fi, err := fs.Stat(filePath)
	require.NoError(t, err)
	require.Greater(t, fi.Size(), int64(0))
}
