package model

import "gorm.io/gorm"

// Staff represents hospital personnel
// @Description Hastane personel bilgileri
type Staff struct {
	gorm.Model   `swaggerignore:"true"`
	HospitalID   uint   `json:"hospital_id" gorm:"not null" example:"1" binding:"required"`            // Hangi hastane
	PolyclinicID *uint  `json:"polyclinic_id,omitempty" example:"1"`                                   // Bağlı olduğu poliklinik (nullable - güvenlik gibi genel personel için)
	FirstName    string `json:"first_name" gorm:"not null" example:"Dr. Mehmet" binding:"required"`    // Ad
	LastName     string `json:"last_name" gorm:"not null" example:"Özkan" binding:"required"`          // Soyad
	TCKN         string `json:"tc" gorm:"unique;not null" example:"98765432101" binding:"required"`    // TC Kimlik No
	Phone        string `json:"phone" gorm:"unique;not null" example:"05559876543" binding:"required"` // Telefon
	JobGroupID   uint   `json:"job_group_id" gorm:"not null" example:"1" binding:"required"`           // Meslek grubu
	JobTitleID   uint   `json:"job_title_id" gorm:"not null" example:"1" binding:"required"`           // Unvan
	WorkDays     string `json:"work_days" gorm:"type:json" example:"[1,2,3,4,5]" binding:"required"`   // Çalışma günleri (JSON array: 1=Pazartesi, 7=Pazar)
	IsActive     bool   `json:"is_active" gorm:"default:true" example:"true"`                          // Aktif mi?

	// İlişkiler
	Hospital   Hospital            `json:"hospital,omitempty" gorm:"foreignKey:HospitalID"`
	Polyclinic *HospitalPolyclinic `json:"polyclinic,omitempty" gorm:"foreignKey:PolyclinicID"`
	JobGroup   JobGroup            `json:"job_group,omitempty" gorm:"foreignKey:JobGroupID"`
	JobTitle   JobTitle            `json:"job_title,omitempty" gorm:"foreignKey:JobTitleID"`
}
