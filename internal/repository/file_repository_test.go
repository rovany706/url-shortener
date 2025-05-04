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
		info *ShortenedURLInfo
		ok   bool
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
				info: &ShortenedURLInfo{
					FullURL:   "http://example.com",
					ShortID:   "89dce6a4",
					IsDeleted: false,
					UserID:    1,
				},
				ok: true,
			},
		},
		{
			name:    "missing ID",
			shortID: "11111111",
			want: want{
				info: nil,
				ok:   false,
			},
		},
	}

	ctx := context.Background()
	testStoragePath := "/home/test/storage.json"
	fs := afero.NewMemMapFs()
	err := fs.MkdirAll("/home/test", 0755)
	require.NoError(t, err)
	loadTestData(t, fs, "testdata/test_storage.json", testStoragePath)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repository, err := NewFileRepository(fs, testStoragePath)
			require.NoError(t, err)

			info, ok := repository.GetFullURL(ctx, tt.shortID)

			assert.Equal(t, tt.want.ok, ok)

			if !ok {
				assert.Nil(t, info)
				return
			}

			assert.Equal(t, *tt.want.info, *info)
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
			err := fs.MkdirAll("/home/test", 0755)
			require.NoError(t, err)
			testStoragePath := "/home/test/storage.json"
			loadTestData(t, fs, "testdata/test_storage.json", testStoragePath)
			fi, err := fs.Stat(testStoragePath)
			require.NoError(t, err)
			testDataFileSize := fi.Size()

			repository, err := NewFileRepository(fs, testStoragePath)
			require.NoError(t, err)

			err = repository.SaveEntry(ctx, 1, tt.shortID, tt.fullURL)
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

func TestSaveEntries(t *testing.T) {
	tests := []struct {
		name             string
		newEntries       map[string]string
		wantWriteNewData bool
	}{
		{
			name: "write new data",
			newEntries: map[string]string{
				"1": "https://ya.ru",
				"2": "https://google.com",
			},
			wantWriteNewData: true,
		},
		{
			name: "no new data (existing shortIDs)",
			newEntries: map[string]string{
				"89dce6a4": "https://ya.ru",
				"ec2c0086": "https://google.com",
			},
			wantWriteNewData: false,
		},
	}

	ctx := context.Background()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := afero.NewMemMapFs()
			err := fs.MkdirAll("/home/test", 0755)
			require.NoError(t, err)
			testStoragePath := "/home/test/storage.json"
			loadTestData(t, fs, "testdata/test_storage.json", testStoragePath)
			fi, err := fs.Stat(testStoragePath)
			require.NoError(t, err)
			testDataFileSize := fi.Size()

			repository, err := NewFileRepository(fs, testStoragePath)
			require.NoError(t, err)

			err = repository.SaveEntries(ctx, 1, tt.newEntries)
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

func TestPing(t *testing.T) {
	ctx := context.Background()
	fs := afero.NewMemMapFs()
	err := fs.MkdirAll("/home/test", 0755)
	require.NoError(t, err)
	testStoragePath := "/home/test/storage.json"
	loadTestData(t, fs, "testdata/test_storage.json", testStoragePath)

	repository, err := NewFileRepository(fs, testStoragePath)
	require.NoError(t, err)

	err = repository.Ping(ctx)
	assert.ErrorIs(t, err, ErrPingNotSupported)
}

func TestClose(t *testing.T) {
	fs := afero.NewMemMapFs()
	err := fs.MkdirAll("/home/test", 0755)
	require.NoError(t, err)
	testStoragePath := "/home/test/storage.json"
	loadTestData(t, fs, "testdata/test_storage.json", testStoragePath)

	repository, err := NewFileRepository(fs, testStoragePath)
	require.NoError(t, err)

	err = repository.Close()
	assert.NoError(t, err)
}

func TestGetShortID(t *testing.T) {
	tests := []struct {
		name        string
		fullURL     string
		wantShortID string
	}{
		{
			name:        "existing full URL",
			fullURL:     "http://example.com",
			wantShortID: "89dce6a4",
		},
		{
			name:        "non existing full URL",
			fullURL:     "https://ya.ru",
			wantShortID: "",
		},
	}

	ctx := context.Background()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := afero.NewMemMapFs()
			err := fs.MkdirAll("/home/test", 0755)
			require.NoError(t, err)
			testStoragePath := "/home/test/storage.json"
			loadTestData(t, fs, "testdata/test_storage.json", testStoragePath)

			repository, err := NewFileRepository(fs, testStoragePath)
			require.NoError(t, err)

			shortID, err := repository.GetShortID(ctx, tt.fullURL)
			require.NoError(t, err)
			assert.Equal(t, tt.wantShortID, shortID)
		})
	}
}
