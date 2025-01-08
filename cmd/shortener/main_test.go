package main

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/rovany706/url-shortener/internal/app"
	"github.com/rovany706/url-shortener/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRedirectHandler(t *testing.T) {
	appConfig = &config.AppConfig{BaseURL: "http://localhost:8080", AppRunAddress: ":8080"}

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
			app := app.URLShortenerApp{
				ShortURLMap: tt.shortURLMap,
			}
			request := httptest.NewRequest(http.MethodGet, tt.request, nil)
			w := httptest.NewRecorder()
			// fix chi.URLParams (https://github.com/go-chi/chi/issues/76#issuecomment-370145140)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.requestID)
			request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, rctx))

			RedirectHandler(&app)(w, request)

			result := w.Result()
			defer result.Body.Close()

			assert.Equal(t, tt.want.code, result.StatusCode)
			assert.Equal(t, tt.want.location, result.Header.Get("Location"))
		})
	}
}

func TestMakeShortURLHandler(t *testing.T) {
	appConfig = &config.AppConfig{BaseURL: "http://localhost:8080", AppRunAddress: ":8080"}

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
				body:        "http://localhost:8080/488575e6",
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

			MakeShortURLHandler(&app)(w, request)
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

func testRequest(t *testing.T, ts *httptest.Server, method string, path string, body string) (*http.Response, string) {
	request, err := http.NewRequest(method, ts.URL+path, strings.NewReader(body))
	require.NoError(t, err)

	client := ts.Client()
	// prevent redirects
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	response, err := client.Do(request)
	require.NoError(t, err)

	responseBody, err := io.ReadAll(response.Body)
	require.NoError(t, err)

	return response, string(responseBody)
}

func TestMainRouter(t *testing.T) {
	appConfig = &config.AppConfig{BaseURL: "http://localhost:8080", AppRunAddress: ":8080"}

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
			request:      "/id1",
			method:       http.MethodPut,
			body:         "123",
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
			r := MainRouter(&app)
			ts := httptest.NewServer(r)
			defer ts.Close()

			response, _ := testRequest(t, ts, tt.method, tt.request, tt.body)
			defer response.Body.Close()

			assert.Equal(t, tt.expectedCode, response.StatusCode)
		})
	}
}
