package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// Şifreyi hashle
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// Girilen şifre ile hash aynı mı kontrol et
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
