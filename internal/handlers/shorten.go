package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/rovany706/url-shortener/internal/app"
	"github.com/rovany706/url-shortener/internal/auth"
	"github.com/rovany706/url-shortener/internal/config"
	"github.com/rovany706/url-shortener/internal/models"
	"github.com/rovany706/url-shortener/internal/repository"
	"go.uber.org/zap"
)

func MakeShortURLHandler(app app.URLShortener, authentication auth.JWTAuthentication, repo repository.Repository, appConfig *config.AppConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := getUserIDFromRequest(r.Context(), authentication, repo, r)

		if err != nil {
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		body, err := io.ReadAll(r.Body)

		if err != nil {
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		shortID, err := app.GetShortID(r.Context(), userID, string(body))

		statusCode := http.StatusCreated
		if err != nil {
			if errors.Is(err, repository.ErrConflict) {
				statusCode = http.StatusConflict
			} else {
				http.Error(w, "", http.StatusBadRequest)
				return
			}
		}

		token, err := authentication.CreateToken(userID)
		if err != nil {
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:  auth.AuthCookieName,
			Value: token,
		})
		w.Header().Add("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(statusCode)
		w.Write([]byte(getShortURL(shortID, appConfig)))
	}
}

func MakeShortURLHandlerJSON(app app.URLShortener, appConfig *config.AppConfig, authentication auth.JWTAuthentication, repo repository.Repository, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := getUserIDFromRequest(r.Context(), authentication, repo, r)

		if err != nil {
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		decoder := json.NewDecoder(r.Body)
		var request models.ShortenRequest

		if err := decoder.Decode(&request); err != nil {
			logger.Info("cannot decode request JSON body", zap.Error(err))
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		shortID, err := app.GetShortID(r.Context(), userID, request.URL)

		statusCode := http.StatusCreated
		if err != nil {
			if errors.Is(err, repository.ErrConflict) {
				statusCode = http.StatusConflict
			} else {
				http.Error(w, "", http.StatusBadRequest)
				return
			}
		}
		response := models.ShortenResponse{
			Result: getShortURL(shortID, appConfig),
		}

		token, err := authentication.CreateToken(userID)
		if err != nil {
			logger.Info("error creating token", zap.Error(err))
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:  auth.AuthCookieName,
			Value: token,
		})
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)

		encoder := json.NewEncoder(w)
		if err := encoder.Encode(response); err != nil {
			logger.Info("error encoding response", zap.Error(err))
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}
}

func MakeShortURLBatchHandler(app app.URLShortener, appConfig *config.AppConfig, authentication auth.JWTAuthentication, repository repository.Repository, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := getUserIDFromRequest(r.Context(), authentication, repository, r)

		if err != nil {
			http.Error(w, "", http.StatusBadRequest)
			return
		}

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

		shortIDs, err := app.GetShortIDBatch(r.Context(), userID, fullURLs)

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

		token, err := authentication.CreateToken(userID)
		if err != nil {
			logger.Info("error creating token", zap.Error(err))
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:  auth.AuthCookieName,
			Value: token,
		})
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

func getUserIDFromRequest(ctx context.Context, authentication auth.JWTAuthentication, repository repository.Repository, r *http.Request) (int, error) {
	authCookie, err := r.Cookie(auth.AuthCookieName)

	if err != nil {
		return getNewUserID(ctx, repository)
	}

	claims, err := authentication.GetClaimsFromToken(authCookie.Value)

	if err != nil {
		return getNewUserID(ctx, repository)
	}

	return claims.UserID, nil
}

func getNewUserID(ctx context.Context, repository repository.Repository) (int, error) {
	newUserID, err := repository.GetNewUserID(ctx)

	if err != nil {
		return -1, err
	}

	return newUserID, nil
}
