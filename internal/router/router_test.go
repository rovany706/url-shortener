package router

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rovany706/url-shortener/internal/app"
	"github.com/rovany706/url-shortener/internal/config"
	"github.com/rovany706/url-shortener/internal/repository/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

func gzipCompressString(s string) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	zb := gzip.NewWriter(buf)
	_, err := zb.Write([]byte(s))

	if err != nil {
		return nil, err
	}

	err = zb.Close()

	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func testRequest(t *testing.T, ts *httptest.Server, method string, path string, requestBody string) (*http.Response, string) {
	compressed, err := gzipCompressString(requestBody)
	require.NoError(t, err)

	request, err := http.NewRequest(method, ts.URL+path, bytes.NewReader(compressed))
	request.Header.Set("Accept-Encoding", "gzip")
	request.Header.Set("Content-Encoding", "gzip")

	require.NoError(t, err)

	client := ts.Client()
	// prevent redirects
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	response, err := client.Do(request)
	require.NoError(t, err)

	zr, err := gzip.NewReader(response.Body)
	require.NoError(t, err)

	responseBody, err := io.ReadAll(zr)
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
		{
			name:         "GET /ping test",
			request:      "/ping",
			method:       http.MethodGet,
			body:         "",
			expectedCode: http.StatusOK,
		},
		{
			name:         "POST /ping test",
			request:      "/ping",
			method:       http.MethodPost,
			body:         "",
			expectedCode: http.StatusMethodNotAllowed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			repository := mock.NewMockRepository(ctrl)
			repository.EXPECT().Ping(gomock.Any()).Return(nil).AnyTimes()
			obs, logs := observer.New(zap.InfoLevel)
			logger := zap.New(obs)
			shortener := app.NewMockURLShortener(shortURLMap)

			r := MainRouter(shortener, appConfig, repository, logger)
			ts := httptest.NewServer(r)
			defer ts.Close()

			response, _ := testRequest(t, ts, tt.method, tt.request, tt.body)
			defer response.Body.Close()

			assert.Equal(t, tt.expectedCode, response.StatusCode)
			assert.NotEmpty(t, logs.AllUntimed())
		})
	}
}
