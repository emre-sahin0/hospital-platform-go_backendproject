package model

import "gorm.io/gorm"

// User represents a hospital user
// @Description Hastane kullanıcı bilgileri
type User struct {
	gorm.Model `swaggerignore:"true"`
	FirstName  string `json:"first_name" gorm:"not null" example:"Ahmet" binding:"required"`                            // Ad
	LastName   string `json:"last_name" gorm:"not null" example:"Yılmaz" binding:"required"`                            // Soyad
	TCKN       string `json:"tc" gorm:"unique;not null" example:"12345678901" binding:"required"`                       // Türkiye Cumhuriyeti Kimlik Numarası
	Email      string `json:"email" gorm:"unique;not null" example:"ahmet.yilmaz@example.com" binding:"required,email"` // E-posta adresi
	Phone      string `json:"phone" gorm:"unique;not null" example:"05551234567" binding:"required"`                    // Telefon numarası
	Password   string `json:"password,omitempty" example:"123456" binding:"required,min=6"`                             // Şifre (JSON'dan okuyabilir ama response'da göstermez)
	Role       string `json:"role" example:"worker" binding:"required,oneof=admin worker"`                              // Rol: "admin" veya "worker"
	HospitalID uint   `json:"hospital_id" example:"1" binding:"required"`                                               // Hangi hastaneye ait
}
