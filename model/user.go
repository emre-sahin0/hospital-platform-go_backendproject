package model

import "gorm.io/gorm"

const (
	RoleYetkili = "yetkili" // Yetkililer - tüm işlemleri yapabilir
	RoleCalisan = "çalışan" // Çalışanlar - sadece görüntüleme
)

// @Description Hastane kullanıcı bilgileri
type User struct {
	gorm.Model `swaggerignore:"true"`
	HospitalID uint   `json:"hospital_id" gorm:"not null" example:"1" binding:"required"`                               // Hangi hastaneye ait
	FirstName  string `json:"first_name" gorm:"not null" example:"Ahmet" binding:"required"`                            // Ad
	LastName   string `json:"last_name" gorm:"not null" example:"Yılmaz" binding:"required"`                            // Soyad
	TCKN       string `json:"tc" gorm:"unique;not null" example:"12345678901" binding:"required"`                       // Türkiye Cumhuriyeti Kimlik Numarası
	Email      string `json:"email" gorm:"unique;not null" example:"ahmet.yilmaz@example.com" binding:"required,email"` // E-posta adresi
	Phone      string `json:"phone" gorm:"unique;not null" example:"05551234567" binding:"required"`                    // Telefon numarası
	Password   string `json:"password,omitempty" example:"123456" binding:"required,min=6"`                             // Şifre (JSON'dan okuyabilir ama response'da göstermez)
	Role       string `json:"role" gorm:"default:çalışan" example:"çalışan" binding:"required,oneof=yetkili çalışan"`   // Rol: "yetkili" veya "çalışan"
	CreatedBy  *uint  `json:"created_by,omitempty" example:"1"`                                                         // Kim tarafından eklendi (nullable - ilk user için)
	IsActive   bool   `json:"is_active" gorm:"default:true" example:"true"`                                             // Aktif mi?

	// İlişkiler
	Hospital Hospital `json:"hospital,omitempty" gorm:"foreignKey:HospitalID"` // Hastane bilgisi
	Creator  *User    `json:"creator,omitempty" gorm:"foreignKey:CreatedBy"`   // Oluşturan kullanıcı
}
