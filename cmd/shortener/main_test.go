package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/rovany706/url-shortener/internal/app"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRedirectHandler(t *testing.T) {
	type want struct {
		code     int
		location string
	}
	tests := []struct {
		name        string
		request     string
		shortURLMap map[string]string
		want        want
	}{
		{
			name:    "redirect test",
			request: "/id1",
			shortURLMap: map[string]string{
				"id1": "http://example.com/",
			},
			want: want{
				code:     http.StatusTemporaryRedirect,
				location: "http://example.com/",
			},
		},
		{
			name:    "error test",
			request: "/id2",
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
			app := app.URLShortenerApp{
				ShortURLMap: tt.shortURLMap,
			}
			request := httptest.NewRequest(http.MethodGet, tt.request, nil)
			w := httptest.NewRecorder()
			RedirectHandler(&app, w, request)

			result := w.Result()
			defer result.Body.Close()

			assert.Equal(t, tt.want.code, result.StatusCode)
			assert.Equal(t, tt.want.location, result.Header.Get("Location"))
		})
	}
}

func TestMakeShortURLHandler(t *testing.T) {
	type want struct {
		statusCode  int
		contentType string
		body        string
	}
	tests := []struct {
		name string
		body string
		want want
	}{
		{
			name: "valid url test",
			body: "http://example.com/123",
			want: want{
				statusCode:  http.StatusCreated,
				contentType: "text/plain; charset=utf-8",
				body:        BaseURL + "488575e6",
			},
		},
		{
			name: "invalid url test",
			body: "http,,//example.com/",
			want: want{
				statusCode:  http.StatusBadRequest,
				contentType: "text/plain; charset=utf-8",
				body:        "\n",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := app.URLShortenerApp{}
			request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.body))
			w := httptest.NewRecorder()

			MakeShortURLHandler(&app, w, request)
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

func TestMainHook(t *testing.T) {
	shortURLMap := map[string]string{
		"id1": "http://example.com/",
	}

	tests := []struct {
		name         string
		request      string
		method       string
		body         string
		expectedCode int
	}{
		{
			name:         "POST test",
			request:      "/",
			method:       http.MethodPost,
			body:         "http://example.com/123",
			expectedCode: http.StatusCreated,
		},
		{
			name:         "GET test",
			request:      "/id1",
			method:       http.MethodGet,
			body:         "",
			expectedCode: http.StatusTemporaryRedirect,
		},
		{
			name:         "PUT test",
			request:      "/",
			method:       http.MethodPut,
			body:         "",
			expectedCode: http.StatusMethodNotAllowed,
		},
		{
			name:         "DELETE test",
			request:      "/",
			method:       http.MethodDelete,
			body:         "",
			expectedCode: http.StatusMethodNotAllowed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := app.URLShortenerApp{
				ShortURLMap: shortURLMap,
			}

			request := httptest.NewRequest(tt.method, tt.request, strings.NewReader(tt.body))
			w := httptest.NewRecorder()
			MainHook(&app)(w, request)

			response := w.Result()
			defer response.Body.Close()

			assert.Equal(t, tt.expectedCode, response.StatusCode)
		})
	}
}
