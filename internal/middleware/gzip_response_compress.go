package middleware

import (
	"compress/gzip"
	"net/http"
	"strings"

	"go.uber.org/zap"
)

type compressWriter struct {
	w  http.ResponseWriter
	zw *gzip.Writer
}

func newCompressWriter(w http.ResponseWriter) *compressWriter {
	return &compressWriter{
		w:  w,
		zw: gzip.NewWriter(w),
	}
}

// Header возвращает заголовок запроса
func (c *compressWriter) Header() http.Header {
	return c.w.Header()
}

// WriteHeader отправляет HTTP-код
func (c *compressWriter) WriteHeader(statusCode int) {
	c.w.Header().Set("Content-Encoding", "gzip")
	c.w.WriteHeader(statusCode)
}

// Write записывает сжатые данные
func (c *compressWriter) Write(p []byte) (int, error) {
	return c.zw.Write(p)
}

// Close закрывает writer
func (c *compressWriter) Close() error {
	return c.zw.Close()
}

// ResponseGzipCompress middleware для сжатия ответа
func ResponseGzipCompress(logger *zap.Logger) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			responseWriter := w

			acceptEncoding := r.Header.Get("Accept-Encoding")
			supportsGzip := strings.Contains(acceptEncoding, "gzip")

			if supportsGzip {
				compressWriter := newCompressWriter(w)
				responseWriter = compressWriter

				defer func() {
					dErr := compressWriter.Close()
					if dErr != nil {
						logger.Error("error closing gzip writer", zap.Error(dErr))
						w.WriteHeader(http.StatusInternalServerError)
					}
				}()
			}

			h.ServeHTTP(responseWriter, r)
		}

		return http.HandlerFunc(fn)
	}
}
