package repository

import (
	"hospital-platform/database"
	"hospital-platform/model"
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

// Delete hastaneyi soft delete yapar
func (r *HospitalRepository) Delete(id uint) error {
	result := database.DB.Delete(&model.Hospital{}, id)
	return result.Error
}
