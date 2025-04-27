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

type UserHandlers struct {
	appConfig     *config.AppConfig
	logger        *zap.Logger
	repository    repository.Repository
	tokenManager  auth.TokenManager
	deleteService service.DeleteService
}

func NewUserHandlers(deleteService service.DeleteService, tokenManager auth.TokenManager, repository repository.Repository, appConfig *config.AppConfig, logger *zap.Logger) UserHandlers {
	return UserHandlers{
		appConfig:     appConfig,
		logger:        logger,
		repository:    repository,
		tokenManager:  tokenManager,
		deleteService: deleteService,
	}
}

// GetUserURLsHandler возвращает пользователю список сокращенных им ссылок
func (h *UserHandlers) GetUserURLsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := getUserIDFromRequest(r.Context(), h.tokenManager, h.repository, r)

		if err != nil {
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		if userID < 1 {
			h.logger.Info("user id is invalid")
			http.Error(w, "", http.StatusUnauthorized)
			return
		}

		token, err := h.tokenManager.CreateToken(userID)
		if err != nil {
			h.logger.Info("error creating token", zap.Error(err))
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:  auth.AuthCookieName,
			Value: token,
		})

		shortIDMap, err := h.repository.GetUserEntries(r.Context(), userID)

		if err != nil {
			h.logger.Info("error getting user urls", zap.Error(err))
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
				ShortURL:    getShortURL(shortID, h.appConfig),
				OriginalURL: fullURL,
			})
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		encoder := json.NewEncoder(w)
		if err := encoder.Encode(response); err != nil {
			h.logger.Info("error encoding response", zap.Error(err))
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}
}

// DeleteUserURLsHandler принимает запросы на удаление сокращенных ссылкок
func (h *UserHandlers) DeleteUserURLsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := getUserIDFromRequest(r.Context(), h.tokenManager, h.repository, r)

		if err != nil {
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		decoder := json.NewDecoder(r.Body)
		var request models.DeleteURLsRequest

		if err = decoder.Decode(&request); err != nil {
			h.logger.Info("cannot decode request JSON body", zap.Error(err))
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

		h.deleteService.Put(deleteChan)

		w.WriteHeader(http.StatusAccepted)
	}
}
