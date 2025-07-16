package model

// LoginRequest represents login credentials
// @Description Kullanıcı giriş bilgileri
type LoginRequest struct {
	Email    string `json:"email" example:"ahmet.yilmaz@example.com" binding:"required,email"` // E-posta adresi
	Password string `json:"password" example:"123456" binding:"required,min=6"`                // Şifre
}
