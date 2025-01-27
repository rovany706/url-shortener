package router

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
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

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
	appConfig := config.NewConfig(config.WithBaseURL("http://localhost:8080"), config.WithAppRunAddress(":8080"))

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
			name:         "POST / test",
			request:      "/",
			method:       http.MethodPost,
			body:         "http://example.com/123",
			expectedCode: http.StatusCreated,
		},
		{
			name:         "GET / test",
			request:      "/id1",
			method:       http.MethodGet,
			body:         "",
			expectedCode: http.StatusTemporaryRedirect,
		},
		{
			name:         "POST /api/shorten test",
			request:      "/api/shorten",
			method:       http.MethodPost,
			body:         "{\"result\":\"http://localhost:8080/0\"}\n",
			expectedCode: http.StatusCreated,
		},
		{
			name:         "GET /api/shorten test",
			request:      "/api/shorten",
			method:       http.MethodGet,
			body:         "",
			expectedCode: http.StatusMethodNotAllowed,
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
			obs, logs := observer.New(zap.InfoLevel)
			logger := zap.New(obs)
			shortener := app.NewMockURLShortener(shortURLMap)

			r := MainRouter(shortener, appConfig, logger)
			ts := httptest.NewServer(r)
			defer ts.Close()

			response, _ := testRequest(t, ts, tt.method, tt.request, tt.body)
			defer response.Body.Close()

			assert.Equal(t, tt.expectedCode, response.StatusCode)
			assert.NotEmpty(t, logs.AllUntimed())
		})
	}
}
