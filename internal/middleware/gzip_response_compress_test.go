package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestResponseGzipCompress(t *testing.T) {
	responseBody := `{"result":"http://localhost:8080/0"}`

	request := httptest.NewRequest(http.MethodGet, "/", nil)
	request.Header.Set("Accept-Encoding", "gzip")

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(responseBody))
		w.WriteHeader(http.StatusOK)
	})

	middlewareHandler := ResponseGzipCompress()(handler)

	recorder := httptest.NewRecorder()
	middlewareHandler.ServeHTTP(recorder, request)

	resp := recorder.Result()
	defer resp.Body.Close()

	zr, err := gzip.NewReader(resp.Body)
	require.NoError(t, err)

	b, err := io.ReadAll(zr)
	require.NoError(t, err)

	require.JSONEq(t, responseBody, string(b))
}
