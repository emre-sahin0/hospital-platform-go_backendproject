package repository

import (
	"hospital-platform/database"
	"hospital-platform/model"
)

// LocationRepository il ve ilçe veritabanı işlemlerini yönetir
type LocationRepository struct{}

// NewLocationRepository yeni bir lokasyon repository'si oluşturur
func NewLocationRepository() *LocationRepository {
	return &LocationRepository{}
}

// GetAllProvinces tüm illeri getirir
func (r *LocationRepository) GetAllProvinces() ([]model.Province, error) {
	var provinces []model.Province
	result := database.DB.Order("name ASC").Find(&provinces)
	return provinces, result.Error
}

// GetProvinceByID ID'ye göre il bilgilerini getirir
func (r *LocationRepository) GetProvinceByID(id uint) (*model.Province, error) {
	var province model.Province
	result := database.DB.First(&province, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &province, nil
}

// GetDistrictsByProvinceID belirtilen ile ait tüm ilçeleri getirir
func (r *LocationRepository) GetDistrictsByProvinceID(provinceID uint) ([]model.District, error) {
	var districts []model.District
	result := database.DB.Where("province_id = ?", provinceID).Order("name ASC").Find(&districts)
	return districts, result.Error
}

// GetDistrictByID ID'ye göre ilçe bilgilerini getirir
func (r *LocationRepository) GetDistrictByID(id uint) (*model.District, error) {
	var district model.District
	result := database.DB.Preload("Province").First(&district, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &district, nil
}

// ValidateProvinceDistrict ilçenin belirtilen ile ait olup olmadığını kontrol eder
func (r *LocationRepository) ValidateProvinceDistrict(provinceID, districtID uint) (bool, error) {
	var count int64
	result := database.DB.Model(&model.District{}).
		Where("id = ? AND province_id = ?", districtID, provinceID).
		Count(&count)

	if result.Error != nil {
		return false, result.Error
	}

	return count > 0, nil
}
