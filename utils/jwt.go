package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var JWTSecret = []byte("supersecretkey") // Production'da config'den alınmalı

type Claims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"` // Yetkili veya Çalışan
	jwt.RegisteredClaims
}

// GenerateJWT kullanıcı için JWT token oluşturur
func GenerateJWT(userID uint, email, role string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour) // 24 saat geçerli

	claims := &Claims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JWTSecret)
}

// ValidateJWT JWT token'ını doğrular ve claims'leri döndürür
func ValidateJWT(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return JWTSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}

	return claims, nil
}
