package model

import "gorm.io/gorm"

// Polyclinic represents a hospital polyclinic
// @Description Hastane poliklinik bilgileri
type Polyclinic struct {
	gorm.Model `swaggerignore:"true"`
	Name       string `json:"name" gorm:"not null" example:"Kardiyoloji" binding:"required"` // Poliklinik adı
	Floor      int    `json:"floor" gorm:"not null" example:"2" binding:"required"`          // Kat numarası
	RoomNumber int    `json:"room_number" gorm:"not null" example:"205" binding:"required"`  // Oda numarası
}
