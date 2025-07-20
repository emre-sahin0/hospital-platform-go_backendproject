package repository

import (
	"fmt"
	"hospital-platform/database"
	"hospital-platform/model"
	"hospital-platform/utils"
)

// HospitalRepository hastane veritabanı işlemlerini yönetir
type HospitalRepository struct{}

// NewHospitalRepository yeni bir hastane repository'si oluşturur
func NewHospitalRepository() *HospitalRepository {
	return &HospitalRepository{}
}

// Create veritabanında yeni bir hastane oluşturur
func (r *HospitalRepository) Create(hospital *model.Hospital) error {
	result := database.DB.Create(hospital)
	return result.Error
}

// GetByID ID'ye göre hastane bilgilerini getirir
func (r *HospitalRepository) GetByID(id uint) (*model.Hospital, error) {
	var hospital model.Hospital
	result := database.DB.Preload("Province").Preload("District").First(&hospital, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &hospital, nil
}

// GetByTaxID vergi kimlik numarasına göre hastane olup olmadığını kontrol eder
func (r *HospitalRepository) GetByTaxID(taxID string) (*model.Hospital, error) {
	var hospital model.Hospital
	result := database.DB.Where("tax_id = ?", taxID).First(&hospital)
	if result.Error != nil {
		return nil, result.Error
	}
	return &hospital, nil
}

// GetByEmail e-posta adresine göre hastane olup olmadığını kontrol eder
func (r *HospitalRepository) GetByEmail(email string) (*model.Hospital, error) {
	var hospital model.Hospital
	result := database.DB.Where("email = ?", email).First(&hospital)
	if result.Error != nil {
		return nil, result.Error
	}
	return &hospital, nil
}

// GetByPhone telefon numarasına göre hastane olup olmadığını kontrol eder
func (r *HospitalRepository) GetByPhone(phone string) (*model.Hospital, error) {
	var hospital model.Hospital
	result := database.DB.Where("phone = ?", phone).First(&hospital)
	if result.Error != nil {
		return nil, result.Error
	}
	return &hospital, nil
}

// Update hastane bilgilerini günceller
func (r *HospitalRepository) Update(hospital *model.Hospital) error {
	result := database.DB.Save(hospital)
	return result.Error
}

// CreateHospitalWithAdmin hastane ve admin kullanıcıyı transaction ile oluşturur
func (r *HospitalRepository) CreateHospitalWithAdmin(req *model.HospitalRegistrationRequest, hashedPassword string) (*model.HospitalRegistrationResponse, error) {
	// Transaction başlat
	tx := database.DB.Begin()
	if tx.Error != nil {
		return nil, fmt.Errorf("transaction başlatılamadı: %v", tx.Error)
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 1. Hastaneyi oluştur
	hospital := &model.Hospital{
		Name:          req.HospitalName,
		TaxID:         req.TaxID,
		Email:         req.HospitalEmail,
		Phone:         req.HospitalPhone,
		ProvinceID:    req.ProvinceID,
		DistrictID:    req.DistrictID,
		AddressDetail: req.AddressDetail,
	}

	if err := tx.Create(hospital).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("hastane oluşturulamadı: %v", err)
	}

	// 2. Admin kullanıcıyı oluştur
	adminUser := &model.User{
		HospitalID: hospital.ID,
		FirstName:  req.AdminFirstName,
		LastName:   req.AdminLastName,
		TCKN:       req.AdminTCKN,
		Email:      req.AdminEmail,
		Phone:      req.AdminPhone,
		Password:   hashedPassword,
		Role:       model.RoleYetkili, // İlk kullanıcı yetkili
		CreatedBy:  nil,               // İlk kullanıcı için nil
		IsActive:   true,
	}

	if err := tx.Create(adminUser).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("admin kullanıcı oluşturulamadı: %v", err)
	}

	// 3. Transaction'ı commit et
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("transaction commit edilemedi: %v", err)
	}

	// 4. JWT token oluştur
	token, err := utils.GenerateJWT(adminUser.ID, adminUser.Email, adminUser.Role, adminUser.HospitalID, adminUser.Email)
	if err != nil {
		return nil, fmt.Errorf("token oluşturulamadı: %v", err)
	}

	// 5. Response için ilişkileri yükle
	hospitalWithRelations, _ := r.GetByID(hospital.ID)
	adminUser.Password = "" // Şifreyi response'dan kaldır

	response := &model.HospitalRegistrationResponse{
		Message:   "Hastane ve admin kullanıcı başarıyla oluşturuldu",
		Hospital:  *hospitalWithRelations,
		AdminUser: *adminUser,
		Token:     token,
	}

	return response, nil
}

// Delete hastaneyi soft delete yapar
func (r *HospitalRepository) Delete(id uint) error {
	result := database.DB.Delete(&model.Hospital{}, id)
	return result.Error
}
