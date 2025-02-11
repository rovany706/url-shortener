package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rovany706/url-shortener/internal/app"
	"github.com/rovany706/url-shortener/internal/config"
	"github.com/rovany706/url-shortener/internal/models"
	"github.com/rovany706/url-shortener/internal/repository"
	"go.uber.org/zap"
)

func RedirectHandler(app app.URLShortener) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		shortID := chi.URLParam(r, "id")
		if fullURL, ok := app.GetFullURL(r.Context(), shortID); ok {
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

		shortID, err := app.GetShortID(r.Context(), string(body))

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

		shortID, err := app.GetShortID(r.Context(), request.URL)

		if err != nil {
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		response := models.ShortenResponse{
			Result: getShortURL(shortID, appConfig),
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

func PingHandler(repository repository.Repository, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := repository.Ping(r.Context())

		if err != nil {
			logger.Info("unable to ping repository data source", zap.Error(err))
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func MakeShortURLBatchHandler(app app.URLShortener, appConfig *config.AppConfig, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		var request models.BatchShortenRequest

		if err := decoder.Decode(&request); err != nil {
			logger.Info("cannot decode request JSON body", zap.Error(err))
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		fullURLs := make([]string, len(request))
		for i, url := range request {
			fullURLs[i] = url.OriginalURL
		}

		shortIDs, err := app.GetShortIDBatch(r.Context(), fullURLs)

		if err != nil {
			logger.Info("error creating short ids", zap.Error(err))
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		responseEntries := make([]models.BatchShortenResponseEntry, len(request))
		for i, shortID := range shortIDs {
			entry := models.BatchShortenResponseEntry{
				CorrelationID: request[i].CorrelationID,
				ShortURL:      getShortURL(shortID, appConfig),
			}

			responseEntries[i] = entry
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		encoder := json.NewEncoder(w)
		if err := encoder.Encode(responseEntries); err != nil {
			logger.Info("error encoding response", zap.Error(err))
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}
}

func getShortURL(shortID string, appConfig *config.AppConfig) string {
	return appConfig.BaseURL + "/" + shortID
}
