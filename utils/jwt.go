package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var JWTSecret = []byte("supersecretkey") // Production'da config'den alınmalı

type Claims struct {
	UserID     uint   `json:"user_id"`
	Email      string `json:"email"`
	Role       string `json:"role"`        // yetkili veya çalışan
	HospitalID uint   `json:"hospital_id"` // Kullanıcının bağlı olduğu hastane
	Username   string `json:"username"`    // Kullanıcı adı
	jwt.RegisteredClaims
}

// GenerateJWT - Kullanıcı için JWT token oluşturur (hospital_id dahil)
// Artık hospital bilgisi de token'da taşınır
func GenerateJWT(userID uint, email, role string, hospitalID uint, username string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour) // 24 saat geçerli

	claims := &Claims{
		UserID:     userID,
		Email:      email,
		Role:       role,
		HospitalID: hospitalID,
		Username:   username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JWTSecret)
}

// ValidateJWT - JWT token'ını doğrular ve claims'leri map olarak döndürür
// Middleware'in beklediği format için map[string]interface{} döner
func ValidateJWT(tokenString string) (map[string]interface{}, error) {
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

	// Claims'leri map formatına çevir
	claimsMap := map[string]interface{}{
		"user_id":     float64(claims.UserID), // JSON'da number float64 olur
		"email":       claims.Email,
		"role":        claims.Role,
		"hospital_id": float64(claims.HospitalID), // JSON'da number float64 olur
		"username":    claims.Username,
	}

	return claimsMap, nil
}

// ValidateJWTWithClaims - JWT token'ını doğrular ve Claims struct olarak döndürür
// Backward compatibility için eski fonksiyon
func ValidateJWTWithClaims(tokenString string) (*Claims, error) {
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
