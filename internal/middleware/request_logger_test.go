package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

func TestRequestLogger(t *testing.T) {
	type want struct {
		method string
		path   string
	}

	tests := []struct {
		name   string
		method string
		path   string
		want   want
	}{
		{
			name:   "GET /",
			method: http.MethodGet,
			path:   "/",
			want: want{
				method: "HTTP",
				path:   "/",
			},
		},
		{
			name:   "POST /test",
			method: http.MethodPost,
			path:   "/test",
			want: want{
				method: "POST",
				path:   "/test",
			},
		},
	}

	emptyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obs, logs := observer.New(zap.InfoLevel)
			logger := zap.New(obs)

			request := httptest.NewRequest(tt.method, tt.path, nil)
			middlewareHandler := RequestLogger(logger)(emptyHandler)

			middlewareHandler.ServeHTTP(httptest.NewRecorder(), request)

			require.Equal(t, 1, logs.Len())

			logEntry := logs.AllUntimed()[0]

			assert.Equal(t, zap.InfoLevel, logEntry.Level)
			assert.Equal(t, "got incoming HTTP request", logEntry.Message)
			assert.Equal(t, tt.method, logEntry.ContextMap()["method"])
			assert.Equal(t, tt.path, logEntry.ContextMap()["path"])

			// can't check duration value because it's time-dependent
			_, durationTagExists := logEntry.ContextMap()["duration"]
			_, durationIsTimeDuration := logEntry.ContextMap()["duration"].(time.Duration)
			assert.True(t, durationTagExists)
			assert.True(t, durationIsTimeDuration)
		})
	}

}
