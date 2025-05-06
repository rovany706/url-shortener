package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap/zaptest"

	"github.com/rovany706/url-shortener/internal/repository/mock"
)

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
			defer ctrl.Finish()

			repository := mock.NewMockRepository(ctrl)
			repository.EXPECT().Ping(gomock.Any()).Return(tt.pingErr)

			request := httptest.NewRequest(http.MethodGet, "/ping", nil)
			w := httptest.NewRecorder()

			PingHandler(repository, testLogger)(w, request)

			response := w.Result()
			defer func() {
				require.NoError(t, response.Body.Close())
			}()

			assert.Equal(t, tt.wantStatusCode, response.StatusCode)
		})
	}
}
