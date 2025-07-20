package model

// LoginRequest represents login credentials
// @Description Kullanıcı giriş bilgileri (email veya telefon ile)
type LoginRequest struct {
	EmailOrPhone string `json:"email_or_phone" example:"ahmet.yilmaz@example.com" binding:"required"` // E-posta adresi veya telefon numarası
	Password     string `json:"password" example:"123456" binding:"required,min=6"`                   // Şifre
}
