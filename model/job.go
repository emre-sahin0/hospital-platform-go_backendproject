package model

import "gorm.io/gorm"

// JobGroup represents a professional group (Doktor, İdari, Hizmet)
// @Description Meslek grubu bilgileri
type JobGroup struct {
	gorm.Model `swaggerignore:"true"`
	Name       string     `json:"name" gorm:"not null;unique" example:"Doktor" binding:"required"` // Meslek grubu adı
	JobTitles  []JobTitle `json:"job_titles,omitempty" gorm:"foreignKey:JobGroupID"`               // Bu gruba ait unvanlar
}

// JobTitle represents a job title within a professional group
// @Description Unvan bilgileri
type JobTitle struct {
	gorm.Model `swaggerignore:"true"`
	JobGroupID uint     `json:"job_group_id" gorm:"not null" example:"1" binding:"required"`    // Bağlı olduğu meslek grubu
	Name       string   `json:"name" gorm:"not null" example:"Uzman Doktor" binding:"required"` // Unvan adı
	IsUnique   bool     `json:"is_unique" gorm:"default:false" example:"false"`                 // Bu unvandan hastanede sadece 1 tane olabilir mi? (örn: Başhekim)
	JobGroup   JobGroup `json:"job_group,omitempty" gorm:"foreignKey:JobGroupID"`               // Meslek grubu bilgisi
}
