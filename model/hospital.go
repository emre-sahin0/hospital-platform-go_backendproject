package model

import "gorm.io/gorm"

// Hospital represents a hospital/healthcare facility
// @Description Hastane bilgileri
type Hospital struct {
	gorm.Model    `swaggerignore:"true"`
	Name          string `json:"name" gorm:"not null" example:"Acıbadem Hastanesi" binding:"required"`                // Hastane adı
	TaxID         string `json:"tax_id" gorm:"unique;not null" example:"1234567890" binding:"required"`               // Vergi kimlik numarası
	Email         string `json:"email" gorm:"unique;not null" example:"info@acibadem.com" binding:"required,email"`   // E-posta adresi
	Phone         string `json:"phone" gorm:"unique;not null" example:"02121234567" binding:"required"`               // Telefon numarası
	ProvinceID    uint   `json:"province_id" gorm:"not null" example:"1" binding:"required"`                          // İl ID
	DistrictID    uint   `json:"district_id" gorm:"not null" example:"1" binding:"required"`                          // İlçe ID
	AddressDetail string `json:"address_detail" gorm:"not null" example:"Beşiktaş Caddesi No:123" binding:"required"` // Açık adres

	// İlişkiler
	Province Province `json:"province,omitempty" gorm:"foreignKey:ProvinceID"` // İl bilgisi
	District District `json:"district,omitempty" gorm:"foreignKey:DistrictID"` // İlçe bilgisi
	Users    []User   `json:"users,omitempty" gorm:"foreignKey:HospitalID"`    // Hastane kullanıcıları
}
