package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/rovany706/url-shortener/internal/app"
	"github.com/rovany706/url-shortener/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

func TestMakeShortURLHandler(t *testing.T) {
	appConfig := config.NewConfig(config.WithBaseURL("http://localhost:8080"), config.WithAppRunAddress(":8080"))

	type want struct {
		statusCode  int
		contentType string
		body        string
	}
	tests := []struct {
		name    string
		body    string
		wantErr bool
		want    want
	}{
		{
			name:    "valid url test",
			body:    "http://example.com/123",
			wantErr: false,
			want: want{
				statusCode:  http.StatusCreated,
				contentType: "text/plain; charset=utf-8",
				body:        "http://localhost:8080/0",
			},
		},
		{
			name:    "invalid url test",
			body:    "http,,//example.com/",
			wantErr: true,
			want: want{
				statusCode:  http.StatusBadRequest,
				contentType: "text/plain; charset=utf-8",
				body:        "\n",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var shortener app.URLShortener
			if tt.wantErr {
				shortener = &app.ErrMockURLShortener{}
			} else {
				shortener = app.NewMockURLShortener(map[string]string{})
			}

			request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.body))
			w := httptest.NewRecorder()

			MakeShortURLHandler(shortener, appConfig)(w, request)
			response := w.Result()

			defer response.Body.Close()
			responseBody, err := io.ReadAll(response.Body)
			require.NoError(t, err)

			assert.Equal(t, tt.want.statusCode, response.StatusCode)
			assert.Equal(t, tt.want.contentType, response.Header.Get("Content-Type"))
			assert.Equal(t, tt.want.body, string(responseBody))
		})
	}
}

func TestMakeShortURLHandlerJSON(t *testing.T) {
	appConfig := config.NewConfig(config.WithBaseURL("http://localhost:8080"), config.WithAppRunAddress(":8080"))

	type want struct {
		statusCode  int
		contentType string
		body        string
	}
	tests := []struct {
		name    string
		body    string
		wantErr bool
		want    want
	}{
		{
			name:    "valid url and json test",
			body:    `{"url": "http://example.com"}`,
			wantErr: false,
			want: want{
				statusCode:  http.StatusCreated,
				contentType: "application/json",
				body:        "{\"result\":\"http://localhost:8080/0\"}\n",
			},
		},
		{
			name:    "invalid url test",
			body:    `{"url": "http:,,example.com"}`,
			wantErr: true,
			want: want{
				statusCode:  http.StatusBadRequest,
				contentType: "text/plain; charset=utf-8",
				body:        "\n",
			},
		},
		{
			name:    "invalid key name test",
			body:    `{"URL": "http://example.com"}`,
			wantErr: true,
			want: want{
				statusCode:  http.StatusBadRequest,
				contentType: "text/plain; charset=utf-8",
				body:        "\n",
			},
		},
		{
			name:    "wrong request json model test",
			body:    `{"url": "http://example.com", "alice": "bob"}`,
			wantErr: true,
			want: want{
				statusCode:  http.StatusBadRequest,
				contentType: "text/plain; charset=utf-8",
				body:        "\n",
			},
		},
		{
			name:    "not json test",
			body:    "http://example.com",
			wantErr: true,
			want: want{
				statusCode:  http.StatusBadRequest,
				contentType: "text/plain; charset=utf-8",
				body:        "\n",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testLogger := zaptest.NewLogger(t)
			var shortener app.URLShortener
			if tt.wantErr {
				shortener = &app.ErrMockURLShortener{}
			} else {
				shortener = app.NewMockURLShortener(map[string]string{})
			}

			request := httptest.NewRequest(http.MethodPost, "/api/shorten", strings.NewReader(tt.body))
			w := httptest.NewRecorder()

			MakeShortURLHandlerJSON(shortener, appConfig, testLogger)(w, request)
			response := w.Result()

			defer response.Body.Close()
			responseBody, err := io.ReadAll(response.Body)
			require.NoError(t, err)

			assert.Equal(t, tt.want.statusCode, response.StatusCode)
			assert.Equal(t, tt.want.contentType, response.Header.Get("Content-Type"))
			assert.Equal(t, tt.want.body, string(responseBody))
		})
	}
}
