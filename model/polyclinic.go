package model

import "gorm.io/gorm"

// @Description Poliklinik türü bilgileri (master data)
type PolyclinicType struct {
	gorm.Model  `swaggerignore:"true"`
	Name        string `json:"name" gorm:"not null;unique" example:"Kardiyoloji" binding:"required"` // Poliklinik türü adı
	Description string `json:"description" example:"Kalp ve damar hastalıkları"`                     // Açıklama
}

// @Description Hastane poliklinik bilgileri
type HospitalPolyclinic struct {
	gorm.Model       `swaggerignore:"true"`
	HospitalID       uint `json:"hospital_id" gorm:"not null" example:"1" binding:"required"`        // Hangi hastane
	PolyclinicTypeID uint `json:"polyclinic_type_id" gorm:"not null" example:"1" binding:"required"` // Poliklinik türü
	Floor            int  `json:"floor" gorm:"not null" example:"2" binding:"required"`              // Kat numarası
	RoomNumber       int  `json:"room_number" gorm:"not null" example:"205" binding:"required"`      // Oda numarası
	IsActive         bool `json:"is_active" gorm:"default:true" example:"true"`                      // Aktif mi?

	// İlişkiler
	Hospital       Hospital       `json:"hospital,omitempty" gorm:"foreignKey:HospitalID"`
	PolyclinicType PolyclinicType `json:"polyclinic_type,omitempty" gorm:"foreignKey:PolyclinicTypeID"`
	Staff          []Staff        `json:"staff,omitempty" gorm:"foreignKey:PolyclinicID"` // Bu poliklinikte çalışan personel
}

// Legacy Polyclinic struct - backward compatibility için kalsın şimdilik

// @Description Hastane poliklinik bilgileri (eski model)
type Polyclinic struct {
	gorm.Model `swaggerignore:"true"`
	Name       string `json:"name" gorm:"not null" example:"Kardiyoloji" binding:"required"` // Poliklinik adı
	Floor      int    `json:"floor" gorm:"not null" example:"2" binding:"required"`          // Kat numarası
	RoomNumber int    `json:"room_number" gorm:"not null" example:"205" binding:"required"`  // Oda numarası
}
