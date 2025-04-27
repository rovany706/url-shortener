package auth

import "github.com/golang-jwt/jwt/v5"

// Claims хранит полезную нагрузку JWT-токена
type Claims struct {
	jwt.RegisteredClaims
	UserID int
}
