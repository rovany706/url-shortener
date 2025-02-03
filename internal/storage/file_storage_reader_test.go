package storage

import (
	"os"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func loadTestData(t *testing.T, fs afero.Fs, testDataFilePath string, mockFilePath string) {
	testData, err := os.ReadFile(testDataFilePath)
	require.NoError(t, err)

	err = afero.WriteFile(fs, mockFilePath, testData, 0644)
	require.NoError(t, err)
}

func TestNewFileStorageReader(t *testing.T) {
	fs := afero.NewMemMapFs()
	filePath := "home/test/storage.json"
	_, err := NewFileStorageReader(fs, filePath)
	require.NoError(t, err)

	stat, err := fs.Stat(filePath)
	require.NoError(t, err)
	require.NotEmpty(t, stat)
}

func TestReadEntry(t *testing.T) {
	want := StorageEntry{
		ShortID: "89dce6a4",
		FullURL: "http://example.com",
	}

	fs := afero.NewMemMapFs()
	testStoragePath := "home/test/storage.json"
	fs.MkdirAll("/home/test", 0755)
	loadTestData(t, fs, "testdata/test_storage.json", testStoragePath)

	reader, err := NewFileStorageReader(fs, testStoragePath)
	require.NoError(t, err)

	entry, err := reader.ReadEntry()

	require.NoError(t, err)
	assert.Equal(t, want, *entry)
}

func TestReadAllEntries(t *testing.T) {
	want := []StorageEntry{
		{
			ShortID: "89dce6a4",
			FullURL: "http://example.com",
		},
		{
			ShortID: "ec2c0086",
			FullURL: "https://practicum.yandex.ru",
		},
	}

	fs := afero.NewMemMapFs()
	testStoragePath := "home/test/storage.json"
	fs.MkdirAll("/home/test", 0755)
	loadTestData(t, fs, "testdata/test_storage.json", testStoragePath)

	reader, err := NewFileStorageReader(fs, testStoragePath)
	require.NoError(t, err)

	entries, err := reader.ReadAllEntries()

	require.NoError(t, err)
	assert.ElementsMatch(t, want, entries)
}
