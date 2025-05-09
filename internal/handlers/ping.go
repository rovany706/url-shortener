package handlers

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/rovany706/url-shortener/internal/repository"
)

// PingHandler хэндлер для проверки подключения к БД
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
