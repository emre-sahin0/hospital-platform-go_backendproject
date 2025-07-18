package service

import (
	"hospital-platform/model"
)

// LocationService - Coğrafi lokasyon verilerini yöneten servis
// Cache service ile performanslı veri erişimi sağlar
type LocationService struct {
	cacheService *CacheService // Redis cache için
}

// NewLocationService - Yeni bir lokasyon servisi oluşturur
// Cache service'i dependency injection ile alır
func NewLocationService() *LocationService {
	return &LocationService{
		cacheService: NewCacheService(),
	}
}

// GetAllProvinces - Dropdown için tüm illeri cache'den getirir
// İlk çağrıda database'den yükler, sonraki çağrılarda cache'den döner
func (s *LocationService) GetAllProvinces() ([]model.Province, error) {
	return s.cacheService.GetProvinces()
}

// GetDistrictsByProvinceID - Belirtilen ile ait ilçeleri cache'den getirir
// İl ID'ye göre cache anahtarı oluşturur ve cache'den döner
func (s *LocationService) GetDistrictsByProvinceID(provinceID uint) ([]model.District, error) {
	return s.cacheService.GetDistrictsByProvinceID(provinceID)
}
