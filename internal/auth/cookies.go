package auth

import (
	"net/http"

	"go.uber.org/zap"
)

// SetAuthCookie создает JWT-токен и записывает его в виде cookie
func SetAuthCookie(tokenManager TokenManager, w http.ResponseWriter, userID int, logger *zap.Logger) error {
	token, err := tokenManager.CreateToken(userID)
	if err != nil {
		logger.Info("error creating token", zap.Error(err))
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name:  AuthCookieName,
		Value: token,
	})

	return nil
}
