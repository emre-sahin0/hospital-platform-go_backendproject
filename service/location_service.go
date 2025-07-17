package service

import (
	"hospital-platform/model"
	"hospital-platform/repository"
)

// LocationService lokasyon iş mantığını yönetir
type LocationService struct {
	locationRepo *repository.LocationRepository
}

// NewLocationService yeni bir lokasyon servisi oluşturur
func NewLocationService() *LocationService {
	return &LocationService{
		locationRepo: repository.NewLocationRepository(),
	}
}

// GetAllProvinces dropdown için tüm illeri getirir
func (s *LocationService) GetAllProvinces() ([]model.Province, error) {
	return s.locationRepo.GetAllProvinces()
}

// GetDistrictsByProvinceID belirtilen ile ait ilçeleri getirir
func (s *LocationService) GetDistrictsByProvinceID(provinceID uint) ([]model.District, error) {
	return s.locationRepo.GetDistrictsByProvinceID(provinceID)
}
