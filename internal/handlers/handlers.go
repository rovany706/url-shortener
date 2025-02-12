package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rovany706/url-shortener/internal/app"
	"github.com/rovany706/url-shortener/internal/config"
	"github.com/rovany706/url-shortener/internal/models"
	"go.uber.org/zap"
)

func RedirectHandler(app app.URLShortener) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		shortID := chi.URLParam(r, "id")
		if fullURL, ok := app.GetFullURL(shortID); ok {
			http.Redirect(w, r, fullURL, http.StatusTemporaryRedirect)
		} else {
			http.Error(w, "400 Bad Request", http.StatusBadRequest)
		}
	}
}

func MakeShortURLHandler(app app.URLShortener, appConfig *config.AppConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)

		if err != nil {
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		shortID, err := app.GetShortID(string(body))

		if err != nil {
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		w.Header().Add("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(appConfig.BaseURL + "/" + shortID))
	}
}

func MakeShortURLHandlerJSON(app app.URLShortener, appConfig *config.AppConfig, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		var request models.ShortenRequest

		if err := decoder.Decode(&request); err != nil {
			logger.Info("cannot decode request JSON body", zap.Error(err))
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		shortID, err := app.GetShortID(request.URL)

		if err != nil {
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		response := models.ShortenResponse{
			Result: appConfig.BaseURL + "/" + shortID,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		encoder := json.NewEncoder(w)
		if err := encoder.Encode(response); err != nil {
			logger.Info("error encoding response", zap.Error(err))
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}
}
