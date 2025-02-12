package repository

import (
	"context"
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

func TestGetFullURL(t *testing.T) {
	type want struct {
		fullURL string
		ok      bool
	}
	tests := []struct {
		name    string
		shortID string
		want    want
	}{
		{
			name:    "existing ID",
			shortID: "89dce6a4",
			want: want{
				fullURL: "http://example.com",
				ok:      true,
			},
		},
		{
			name:    "missing ID",
			shortID: "11111111",
			want: want{
				fullURL: "",
				ok:      false,
			},
		},
	}

	ctx := context.Background()
	testStoragePath := "/home/test/storage.json"
	fs := afero.NewMemMapFs()
	fs.MkdirAll("/home/test", 0755)
	loadTestData(t, fs, "testdata/test_storage.json", testStoragePath)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repository, err := NewFileRepository(fs, testStoragePath)
			require.NoError(t, err)

			fullURL, ok := repository.GetFullURL(ctx, tt.shortID)

			assert.Equal(t, tt.want.fullURL, fullURL)
			assert.Equal(t, tt.want.ok, ok)
		})
	}
}

func TestSaveEntry(t *testing.T) {
	tests := []struct {
		name             string
		shortID          string
		fullURL          string
		wantWriteNewData bool
	}{
		{
			name:             "write new data",
			shortID:          "1",
			fullURL:          "https://ya.ru",
			wantWriteNewData: true,
		},
		{
			name:             "no new data (existing shortID)",
			shortID:          "89dce6a4",
			fullURL:          "https://ya.ru",
			wantWriteNewData: false,
		},
	}

	ctx := context.Background()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := afero.NewMemMapFs()
			fs.MkdirAll("/home/test", 0755)
			testStoragePath := "/home/test/storage.json"
			loadTestData(t, fs, "testdata/test_storage.json", testStoragePath)
			fi, err := fs.Stat(testStoragePath)
			require.NoError(t, err)
			testDataFileSize := fi.Size()

			repository, err := NewFileRepository(fs, testStoragePath)
			require.NoError(t, err)

			err = repository.SaveEntry(ctx, tt.shortID, tt.fullURL)
			require.NoError(t, err)

			fi, err = fs.Stat(testStoragePath)
			require.NoError(t, err)
			newFileSize := fi.Size()

			if tt.wantWriteNewData {
				assert.Less(t, testDataFileSize, newFileSize)
				return
			}

			assert.Equal(t, testDataFileSize, newFileSize)
		})
	}
}
