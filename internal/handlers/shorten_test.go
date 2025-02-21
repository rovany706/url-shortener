package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/rovany706/url-shortener/internal/app"
	"github.com/stretchr/testify/assert"
)

func TestRedirectHandler(t *testing.T) {
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
			shortener := app.NewMockURLShortener(tt.shortURLMap)

			request := httptest.NewRequest(http.MethodGet, tt.request, nil)
			w := httptest.NewRecorder()
			// fix chi.URLParams (https://github.com/go-chi/chi/issues/76#issuecomment-370145140)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.requestID)
			request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, rctx))

			RedirectHandler(shortener)(w, request)

			result := w.Result()
			defer result.Body.Close()

			assert.Equal(t, tt.want.code, result.StatusCode)
			assert.Equal(t, tt.want.location, result.Header.Get("Location"))
		})
	}
}
