package middleware

import (
	"net/http"

	"github.com/rovany706/url-shortener/internal/auth"
	"go.uber.org/zap"
)

func JWTAuthMiddleware(authentication auth.JWTAuthentication, logger *zap.Logger) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		authFn := func(w http.ResponseWriter, r *http.Request) {
			authCookie, err := r.Cookie(auth.AuthCookieName)

			if err != nil {
				logger.Info("token cookie not found", zap.Error(err))
				http.Error(w, "", http.StatusUnauthorized)
				return
			}

			_, err = authentication.GetClaimsFromToken(authCookie.Value)

			if err != nil {
				logger.Info("token cookie is invalid", zap.Error(err))
			}

			h.ServeHTTP(w, r)
		}

		return http.HandlerFunc(authFn)
	}
}
