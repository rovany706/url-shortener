package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/rovany706/url-shortener/internal/auth"
	"github.com/rovany706/url-shortener/internal/config"
	"github.com/rovany706/url-shortener/internal/models"
	"github.com/rovany706/url-shortener/internal/repository"
	"github.com/rovany706/url-shortener/internal/service"
	"go.uber.org/zap"
)

func GetUserURLs(authentication auth.JWTAuthentication, repository repository.Repository, appConfig *config.AppConfig, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := getUserIDFromRequest(r.Context(), authentication, repository, r)

		if err != nil {
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		if userID < 1 {
			logger.Info("user id is invalid")
			http.Error(w, "", http.StatusUnauthorized)
			return
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

		shortIDMap, err := repository.GetUserEntries(r.Context(), userID)

		if err != nil {
			logger.Info("error getting user urls", zap.Error(err))
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		if len(shortIDMap) == 0 {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		response := make(models.UserShortenedURLs, 0)

		for shortID, fullURL := range shortIDMap {
			response = append(response, models.UserShortenedURL{
				ShortURL:    getShortURL(shortID, appConfig),
				OriginalURL: fullURL,
			})
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		encoder := json.NewEncoder(w)
		if err := encoder.Encode(response); err != nil {
			logger.Info("error encoding response", zap.Error(err))
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}
}

func DeleteUserURLs(deleteService *service.DeleteService, authentication auth.JWTAuthentication, repository repository.Repository, appConfig *config.AppConfig, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := getUserIDFromRequest(r.Context(), authentication, repository, r)

		if err != nil {
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		decoder := json.NewDecoder(r.Body)
		var request models.DeleteURLsRequest

		if err = decoder.Decode(&request); err != nil {
			logger.Info("cannot decode request JSON body", zap.Error(err))
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		deleteChan := make(chan models.UserDeleteRequest)
		go func() {
			defer close(deleteChan)
			for _, shortID := range request {
				deleteRequest := models.UserDeleteRequest{
					UserID:          userID,
					ShortIDToDelete: shortID,
				}

				deleteChan <- deleteRequest
			}
		}()

		deleteService.Put(deleteChan)

		w.WriteHeader(http.StatusAccepted)
	}
}
