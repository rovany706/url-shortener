package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

func TestResponseLogger(t *testing.T) {
	tests := []struct {
		name   string
		body   string
		status int
		want   observer.LoggedEntry
	}{
		{
			name:   "200, empty body",
			body:   "",
			status: http.StatusOK,
			want: observer.LoggedEntry{
				Entry:   zapcore.Entry{Level: zap.InfoLevel, Message: "sent response"},
				Context: []zapcore.Field{zap.Int("size", 0), zap.Int("status", 200)},
			},
		},
		{
			name:   "404, some body (once told me)",
			body:   "the world is gonna roll me",
			status: http.StatusNotFound,
			want: observer.LoggedEntry{
				Entry:   zapcore.Entry{Level: zap.InfoLevel, Message: "sent response"},
				Context: []zapcore.Field{zap.Int("size", len("the world is gonna roll me")), zap.Int("status", 404)},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obs, logs := observer.New(zap.InfoLevel)
			logger := zap.New(obs)

			request := httptest.NewRequest(http.MethodGet, "/", nil)

			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				_, err := w.Write([]byte(tt.body))
				require.NoError(t, err)

				w.WriteHeader(tt.status)
			})

			middlewareHandler := ResponseLogger(logger)(handler)

			middlewareHandler.ServeHTTP(httptest.NewRecorder(), request)

			require.Equal(t, 1, logs.Len())
			assert.Equal(t, tt.want, logs.AllUntimed()[0])
		})
	}
}
