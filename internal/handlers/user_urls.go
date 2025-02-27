package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/rovany706/url-shortener/internal/auth"
	"github.com/rovany706/url-shortener/internal/config"
	"github.com/rovany706/url-shortener/internal/models"
	"github.com/rovany706/url-shortener/internal/repository"
	"go.uber.org/zap"
)

func GetUserURLs(authentication auth.JWTAuthentication, repository repository.Repository, appConfig *config.AppConfig, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authCookie, err := r.Cookie(auth.AuthCookieName)

		if err != nil {
			logger.Info("error finding cookie", zap.Error(err))
			http.Error(w, "", http.StatusUnauthorized)
			return
		}

		claims, err := authentication.GetClaimsFromToken(authCookie.Value)

		if err != nil {
			logger.Info("error parsing claims", zap.Error(err))
			http.Error(w, "", http.StatusUnauthorized)
			return
		}

		shortIDMap, err := repository.GetUserEntries(r.Context(), claims.UserID)

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
