package model

import "gorm.io/gorm"

// Province represents a Turkish province/city
// @Description İl bilgileri
type Province struct {
	gorm.Model `swaggerignore:"true"`
	Name       string     `json:"name" gorm:"not null;unique" example:"İstanbul" binding:"required"` // İl adı
	Districts  []District `json:"districts,omitempty" gorm:"foreignKey:ProvinceID"`                  // İlçeler
}

// District represents a district within a province
// @Description İlçe bilgileri
type District struct {
	gorm.Model `swaggerignore:"true"`
	ProvinceID uint     `json:"province_id" gorm:"not null" example:"1" binding:"required"` // Bağlı olduğu il
	Name       string   `json:"name" gorm:"not null" example:"Beşiktaş" binding:"required"` // İlçe adı
	Province   Province `json:"province,omitempty" gorm:"foreignKey:ProvinceID"`            // İl bilgisi
}
