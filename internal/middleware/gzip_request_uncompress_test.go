package middleware

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

func TestRequestGzipCompress(t *testing.T) {
	testLogger := zaptest.NewLogger(t)
	requestBody := `{"url": "http://example.com"}`
	responseBody := `{"result":"http://localhost:8080/0"}`

	// prepare gzip body
	buf := bytes.NewBuffer(nil)
	zb := gzip.NewWriter(buf)
	_, err := zb.Write([]byte(requestBody))
	require.NoError(t, err)
	err = zb.Close()
	require.NoError(t, err)

	request := httptest.NewRequest(http.MethodPost, "/", buf)
	request.Header.Set("Content-Encoding", "gzip")
	request.Header.Set("Accept-Encoding", "")

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err = w.Write([]byte(responseBody))
		require.NoError(t, err)

		w.WriteHeader(http.StatusOK)
	})

	middlewareHandler := RequestGzipCompress(testLogger)(handler)

	recorder := httptest.NewRecorder()
	middlewareHandler.ServeHTTP(recorder, request)

	resp := recorder.Result()
	defer func() {
		require.NoError(t, resp.Body.Close())
	}()

	b, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.JSONEq(t, responseBody, string(b))
}
