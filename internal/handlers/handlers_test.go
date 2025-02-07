package handlers

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/rovany706/url-shortener/internal/app"
	"github.com/rovany706/url-shortener/internal/config"
	"github.com/rovany706/url-shortener/internal/database/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap/zaptest"
)

func TestRedirectHandler(t *testing.T) {
	type want struct {
		code     int
		location string
	}
	tests := []struct {
		name        string
		requestID   string
		request     string
		shortURLMap map[string]string
		want        want
	}{
		{
			name:      "redirect test",
			requestID: "id1",
			request:   "/id1",
			shortURLMap: map[string]string{
				"id1": "http://example.com/",
			},
			want: want{
				code:     http.StatusTemporaryRedirect,
				location: "http://example.com/",
			},
		},
		{
			name:      "error test",
			requestID: "id2",
			request:   "/id2",
			shortURLMap: map[string]string{
				"id1": "http://example.com/",
			},
			want: want{
				code:     http.StatusBadRequest,
				location: "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shortener := app.NewMockURLShortener(tt.shortURLMap)

			request := httptest.NewRequest(http.MethodGet, tt.request, nil)
			w := httptest.NewRecorder()
			// fix chi.URLParams (https://github.com/go-chi/chi/issues/76#issuecomment-370145140)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.requestID)
			request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, rctx))

			RedirectHandler(shortener)(w, request)

			result := w.Result()
			defer result.Body.Close()

			assert.Equal(t, tt.want.code, result.StatusCode)
			assert.Equal(t, tt.want.location, result.Header.Get("Location"))
		})
	}
}

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

func TestPingHandler(t *testing.T) {
	tests := []struct {
		name           string
		pingErr        error
		wantStatusCode int
	}{
		{
			name:           "ping_successful",
			pingErr:        nil,
			wantStatusCode: http.StatusOK,
		},
		{
			name:           "ping_unsuccessful",
			pingErr:        errors.New("fail"),
			wantStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testLogger := zaptest.NewLogger(t)
			ctrl := gomock.NewController(t)
			db := mock.NewMockDatabase(ctrl)
			db.EXPECT().Ping(gomock.Any()).Return(tt.pingErr)

			request := httptest.NewRequest(http.MethodGet, "/ping", nil)
			w := httptest.NewRecorder()

			PingHandler(db, testLogger)(w, request)

			response := w.Result()
			defer response.Body.Close()

			assert.Equal(t, tt.wantStatusCode, response.StatusCode)
		})
	}
}
