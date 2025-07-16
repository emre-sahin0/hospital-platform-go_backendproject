package model

// ResetPasswordRequest represents a password reset request
// @Description Şifre sıfırlama talebi
type ResetPasswordRequest struct {
	Phone string `json:"phone" example:"05551234567" binding:"required"` // Telefon numarası
}

// ResetPasswordConfirm represents password reset confirmation
// @Description Şifre sıfırlama onayı
type ResetPasswordConfirm struct {
	Phone           string `json:"phone" example:"05551234567" binding:"required"`             // Telefon numarası
	Code            string `json:"code" example:"123456" binding:"required"`                   // Doğrulama kodu
	NewPassword     string `json:"new_password" example:"yeni123456" binding:"required,min=6"` // Yeni şifre
	ConfirmPassword string `json:"confirm_password" example:"yeni123456" binding:"required"`   // Şifre tekrarı
}
