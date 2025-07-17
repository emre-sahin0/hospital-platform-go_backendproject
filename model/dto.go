package model

// HospitalRegistrationRequest represents hospital registration data
// @Description Hastane kayıt verisi
type HospitalRegistrationRequest struct {
	// Hastane Bilgileri
	HospitalName  string `json:"hospital_name" example:"Acıbadem Hastanesi" binding:"required"`
	TaxID         string `json:"tax_id" example:"1234567890" binding:"required"`
	HospitalEmail string `json:"hospital_email" example:"info@acibadem.com" binding:"required,email"`
	HospitalPhone string `json:"hospital_phone" example:"02121234567" binding:"required"`

	// Adres Bilgileri
	ProvinceID    uint   `json:"province_id" example:"1" binding:"required"`
	DistrictID    uint   `json:"district_id" example:"1" binding:"required"`
	AddressDetail string `json:"address_detail" example:"Beşiktaş Caddesi No:123" binding:"required"`

	// Yetkili Bilgileri
	AdminFirstName string `json:"admin_first_name" example:"Ahmet" binding:"required"`
	AdminLastName  string `json:"admin_last_name" example:"Yılmaz" binding:"required"`
	AdminTCKN      string `json:"admin_tc" example:"12345678901" binding:"required"`
	AdminEmail     string `json:"admin_email" example:"ahmet.yilmaz@acibadem.com" binding:"required,email"`
	AdminPhone     string `json:"admin_phone" example:"05551234567" binding:"required"`
	AdminPassword  string `json:"admin_password" example:"123456" binding:"required,min=6"`
}

// HospitalRegistrationResponse represents successful registration response
// @Description Hastane kayıt yanıtı
type HospitalRegistrationResponse struct {
	Message   string   `json:"message" example:"Hastane başarıyla kaydedildi"`
	Hospital  Hospital `json:"hospital"`
	AdminUser User     `json:"admin_user"`
	Token     string   `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

// ValidationError represents validation error details
// @Description Validation hatası detayları
type ValidationError struct {
	Field   string `json:"field" example:"tax_id"`
	Message string `json:"message" example:"Bu vergi kimlik numarası zaten kullanılıyor"`
}

// ==================== POLYCLİNİC DTO'ları ====================

// AddPolyclinicRequest represents adding polyclinic to hospital request
// @Description Hastaneye poliklinik ekleme verisi
type AddPolyclinicRequest struct {
	PolyclinicTypeID uint `json:"polyclinic_type_id" example:"1" binding:"required"` // Master data'dan seçilen tür
	Floor            int  `json:"floor" example:"2" binding:"required"`              // Kat numarası
	RoomNumber       int  `json:"room_number" example:"205" binding:"required"`      // Oda numarası
}

// UpdatePolyclinicRequest represents updating hospital polyclinic request
// @Description Hastane poliklinik güncelleme verisi
type UpdatePolyclinicRequest struct {
	Floor      int  `json:"floor" example:"3" binding:"required"`         // Kat numarası
	RoomNumber int  `json:"room_number" example:"301" binding:"required"` // Oda numarası
	IsActive   bool `json:"is_active" example:"true"`                     // Aktif mi?
}

// HospitalPolyclinicSummary represents hospital polyclinic summary with staff count
// @Description Hastane poliklinik özet bilgileri
type HospitalPolyclinicSummary struct {
	ID                 uint                `json:"id" example:"1"`
	PolyclinicTypeName string              `json:"polyclinic_type_name" example:"Kardiyoloji"`
	Floor              int                 `json:"floor" example:"2"`
	RoomNumber         int                 `json:"room_number" example:"205"`
	IsActive           bool                `json:"is_active" example:"true"`
	TotalStaffCount    int                 `json:"total_staff_count" example:"10"` // Toplam personel sayısı
	StaffByJobGroup    []StaffCountByGroup `json:"staff_by_job_group"`             // Meslek grubuna göre personel sayıları
}

// StaffCountByGroup represents staff count grouped by job group
// @Description Meslek grubuna göre personel sayısı
type StaffCountByGroup struct {
	JobGroupName string `json:"job_group_name" example:"Doktor"`
	Count        int    `json:"count" example:"7"`
}

// ==================== STAFF DTO'ları ====================

// CreateStaffRequest - Yeni personel ekleme işlemi için kullanılan veri yapısı
// Hem zorunlu alanlar hem de isteğe bağlı alanlar içerir
// @Description Hastaneye yeni personel eklerken gönderilecek veriler
type CreateStaffRequest struct {
	FirstName    string `json:"first_name" example:"Dr. Mehmet" binding:"required"` // Personelin adı (zorunlu alan)
	LastName     string `json:"last_name" example:"Özkan" binding:"required"`       // Personelin soyadı (zorunlu alan)
	TCKN         string `json:"tc" example:"98765432101" binding:"required"`        // TC Kimlik numarası - sistemde benzersiz olmalı
	Phone        string `json:"phone" example:"05559876543" binding:"required"`     // Telefon numarası - sistemde benzersiz olmalı
	JobGroupID   uint   `json:"job_group_id" example:"1" binding:"required"`        // Hangi meslek grubuna ait (Doktor, Hemşire vb.)
	JobTitleID   uint   `json:"job_title_id" example:"1" binding:"required"`        // Unvanı (Başhekim, Uzman Doktor vb.) - bazıları unique
	PolyclinicID *uint  `json:"polyclinic_id,omitempty" example:"1"`                // Hangi poliklinikte çalışacak (opsiyonel)
	WorkDays     []int  `json:"work_days" example:"[1,2,3,4,5]" binding:"required"` // Hangi günler çalışacak (1:Pzt, 7:Paz)
}

// UpdateStaffRequest represents updating staff request
// @Description Personel güncelleme verisi
type UpdateStaffRequest struct {
	FirstName    string `json:"first_name" example:"Dr. Ahmet" binding:"required"`    // Ad
	LastName     string `json:"last_name" example:"Yılmaz" binding:"required"`        // Soyad
	Phone        string `json:"phone" example:"05551234567" binding:"required"`       // Telefon
	JobGroupID   uint   `json:"job_group_id" example:"2" binding:"required"`          // Meslek grubu ID
	JobTitleID   uint   `json:"job_title_id" example:"3" binding:"required"`          // Unvan ID
	PolyclinicID *uint  `json:"polyclinic_id,omitempty" example:"2"`                  // Poliklinik ID (nullable)
	WorkDays     []int  `json:"work_days" example:"[1,2,3,4,5,6]" binding:"required"` // Çalışma günleri
	IsActive     bool   `json:"is_active" example:"true"`                             // Aktif mi?
}

// StaffListRequest represents staff filtering and pagination request
// @Description Personel listeleme ve filtreleme verisi
type StaffListRequest struct {
	// Pagination
	Page     int `json:"page" example:"1" binding:"min=1"`               // Sayfa numarası (min: 1)
	PageSize int `json:"page_size" example:"10" binding:"min=1,max=100"` // Sayfa başına kayıt (1-100 arası)

	// Filtering (Optional)
	FirstName    string `json:"first_name,omitempty" example:"Mehmet"` // Ad ile filtreleme
	LastName     string `json:"last_name,omitempty" example:"Özkan"`   // Soyad ile filtreleme
	TCKN         string `json:"tc,omitempty" example:"98765432101"`    // TC ile filtreleme
	JobGroupID   *uint  `json:"job_group_id,omitempty" example:"1"`    // Meslek grubu ile filtreleme
	JobTitleID   *uint  `json:"job_title_id,omitempty" example:"2"`    // Unvan ile filtreleme
	PolyclinicID *uint  `json:"polyclinic_id,omitempty" example:"1"`   // Poliklinik ile filtreleme
	IsActive     *bool  `json:"is_active,omitempty" example:"true"`    // Aktiflik durumu ile filtreleme
}

// StaffListResponse represents paginated staff list response
// @Description Sayfalandırılmış personel listesi yanıtı
type StaffListResponse struct {
	Data       []StaffSummary `json:"data"`       // Personel listesi
	Pagination PaginationInfo `json:"pagination"` // Sayfalama bilgileri
}

// StaffSummary represents staff summary information
// @Description Personel özet bilgileri
type StaffSummary struct {
	ID                 uint    `json:"id" example:"1"`                                       // Personel ID
	FirstName          string  `json:"first_name" example:"Dr. Mehmet"`                      // Ad
	LastName           string  `json:"last_name" example:"Özkan"`                            // Soyad
	TCKN               string  `json:"tc" example:"98765432101"`                             // TC Kimlik No
	Phone              string  `json:"phone" example:"05559876543"`                          // Telefon
	JobGroupName       string  `json:"job_group_name" example:"Doktor"`                      // Meslek grubu adı
	JobTitleName       string  `json:"job_title_name" example:"Uzman Doktor"`                // Unvan adı
	PolyclinicTypeName *string `json:"polyclinic_type_name,omitempty" example:"Kardiyoloji"` // Poliklinik adı (nullable)
	WorkDaysText       string  `json:"work_days_text" example:"Pazartesi-Cuma"`              // Çalışma günleri metni
	IsActive           bool    `json:"is_active" example:"true"`                             // Aktif mi?
}

// PaginationInfo represents pagination metadata
// @Description Sayfalama bilgileri
type PaginationInfo struct {
	CurrentPage  int   `json:"current_page" example:"1"`   // Mevcut sayfa
	PageSize     int   `json:"page_size" example:"10"`     // Sayfa başına kayıt
	TotalRecords int64 `json:"total_records" example:"45"` // Toplam kayıt sayısı
	TotalPages   int   `json:"total_pages" example:"5"`    // Toplam sayfa sayısı
	HasNext      bool  `json:"has_next" example:"true"`    // Sonraki sayfa var mı?
	HasPrev      bool  `json:"has_prev" example:"false"`   // Önceki sayfa var mı?
}
