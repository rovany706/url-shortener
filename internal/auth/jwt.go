package auth

import (
	"crypto/rand"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// TokenExpiryTime срок годности токена
const TokenExpiryTime = time.Hour * 24

// AuthCookieName cookie-ключ токена
const AuthCookieName = "token"

// TokenManager интерфейс менеджера токенов аутентификации
type TokenManager interface {
	GetClaimsFromToken(tokenString string) (*Claims, error)
	CreateToken(userID int) (string, error)
}

// JWTTokenManager реализует TokenManager и использует для работы JWT-токены
type JWTTokenManager struct {
	secretKey []byte
}

// NewJWTTokenManager создает экземпляр JWTTokenManager.
// secretKey - последовательность байт (ключ), используемая для подписи токенов,
// может быть nil для генерации случайного ключа.
func NewJWTTokenManager(secretKey []byte) (*JWTTokenManager, error) {
	var err error
	if secretKey == nil {
		secretKey, err = generateSecretKey()

		if err != nil {
			return nil, err
		}
	}

	return &JWTTokenManager{
		secretKey: secretKey,
	}, nil
}

func generateSecretKey() ([]byte, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

// CreateToken создает JWT-токен для пользователя с userID
func (auth *JWTTokenManager) CreateToken(userID int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExpiryTime)),
		},
		UserID: userID,
	})

	tokenString, err := token.SignedString(auth.secretKey)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// GetClaimsFromToken читает и валидирует JWT-токен.
// Возвращает полезную нагрузку токена.
func (auth *JWTTokenManager) GetClaimsFromToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return auth.secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("token is not valid")
	}

	return claims, nil
}
