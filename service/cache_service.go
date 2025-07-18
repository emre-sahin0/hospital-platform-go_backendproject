package service

import (
	"context"
	"encoding/json"
	"fmt"
	"hospital-platform/database"
	"hospital-platform/model"
	"time"

	"github.com/redis/go-redis/v9" // v9 kullan v8 yerine
)

// CacheService - Redis ile cache işlemlerini yöneten servis
// Master data'lar için cache stratejileri ve TTL yönetimi sağlar
type CacheService struct {
	redisClient *redis.Client
	ctx         context.Context
	defaultTTL  time.Duration // Varsayılan cache süresi
}

// NewCacheService - Yeni cache service instance'ı oluşturur
// Redis bağlantısını initialize eder ve varsayılan TTL'yi ayarlar
func NewCacheService() *CacheService {
	return &CacheService{
		redisClient: database.RedisClient,
		ctx:         context.Background(),
		defaultTTL:  1 * time.Hour, // 1 saat cache süresi
	}
}

// ==================== CACHE ANAHTARLARI ====================

const (
	CACHE_PROVINCES        = "master_data:provinces"           // İller cache anahtarı
	CACHE_DISTRICTS_PREFIX = "master_data:districts:province:" // İlçeler için prefix
	CACHE_JOB_GROUPS       = "master_data:job_groups"          // Meslek grupları
	CACHE_JOB_TITLES       = "master_data:job_titles:group:"   // Unvanlar için prefix
	CACHE_POLYCLINIC_TYPES = "master_data:polyclinic_types"    // Poliklinik tipleri
)

// ==================== İL/İLÇE CACHE İŞLEMLERİ ====================

// GetProvinces - İlleri cache'den getirir, cache miss durumunda database'den yükler
func (cs *CacheService) GetProvinces() ([]model.Province, error) {
	// Cache'den kontrol et
	cachedData, err := cs.redisClient.Get(cs.ctx, CACHE_PROVINCES).Result()
	if err == nil {
		// Cache hit - JSON'dan parse et
		var provinces []model.Province
		if err := json.Unmarshal([]byte(cachedData), &provinces); err == nil {
			return provinces, nil
		}
	}

	// Cache miss - database'den yükle
	var provinces []model.Province
	result := database.DB.Order("name ASC").Find(&provinces)
	if result.Error != nil {
		return nil, fmt.Errorf("iller yüklenemedi: %v", result.Error)
	}

	// Cache'e kaydet
	cs.cacheProvinces(provinces)

	return provinces, nil
}

// GetDistrictsByProvinceID - İle ait ilçeleri cache'den getirir
func (cs *CacheService) GetDistrictsByProvinceID(provinceID uint) ([]model.District, error) {
	cacheKey := fmt.Sprintf("%s%d", CACHE_DISTRICTS_PREFIX, provinceID)

	// Cache'den kontrol et
	cachedData, err := cs.redisClient.Get(cs.ctx, cacheKey).Result()
	if err == nil {
		// Cache hit - JSON'dan parse et
		var districts []model.District
		if err := json.Unmarshal([]byte(cachedData), &districts); err == nil {
			return districts, nil
		}
	}

	// Cache miss - database'den yükle
	var districts []model.District
	result := database.DB.Where("province_id = ?", provinceID).Order("name ASC").Find(&districts)
	if result.Error != nil {
		return nil, fmt.Errorf("ilçeler yüklenemedi: %v", result.Error)
	}

	// Cache'e kaydet
	cs.cacheDistricts(provinceID, districts)

	return districts, nil
}

// ==================== MESLEK GRUPLARI CACHE İŞLEMLERİ ====================

// GetJobGroups - Meslek gruplarını cache'den getirir
func (cs *CacheService) GetJobGroups() ([]model.JobGroup, error) {
	// Cache'den kontrol et
	cachedData, err := cs.redisClient.Get(cs.ctx, CACHE_JOB_GROUPS).Result()
	if err == nil {
		// Cache hit - JSON'dan parse et
		var jobGroups []model.JobGroup
		if err := json.Unmarshal([]byte(cachedData), &jobGroups); err == nil {
			return jobGroups, nil
		}
	}

	// Cache miss - database'den yükle
	var jobGroups []model.JobGroup
	result := database.DB.Order("name ASC").Find(&jobGroups)
	if result.Error != nil {
		return nil, fmt.Errorf("meslek grupları yüklenemedi: %v", result.Error)
	}

	// Cache'e kaydet
	cs.cacheJobGroups(jobGroups)

	return jobGroups, nil
}

// GetJobTitlesByGroupID - Meslek grubuna ait unvanları cache'den getirir
func (cs *CacheService) GetJobTitlesByGroupID(jobGroupID uint) ([]model.JobTitle, error) {
	cacheKey := fmt.Sprintf("%s%d", CACHE_JOB_TITLES, jobGroupID)

	// Cache'den kontrol et
	cachedData, err := cs.redisClient.Get(cs.ctx, cacheKey).Result()
	if err == nil {
		// Cache hit - JSON'dan parse et
		var jobTitles []model.JobTitle
		if err := json.Unmarshal([]byte(cachedData), &jobTitles); err == nil {
			return jobTitles, nil
		}
	}

	// Cache miss - database'den yükle
	var jobTitles []model.JobTitle
	result := database.DB.Where("job_group_id = ?", jobGroupID).Order("name ASC").Find(&jobTitles)
	if result.Error != nil {
		return nil, fmt.Errorf("unvanlar yüklenemedi: %v", result.Error)
	}

	// Cache'e kaydet
	cs.cacheJobTitles(jobGroupID, jobTitles)

	return jobTitles, nil
}

// ==================== POLİKLİNİK TİPLERİ CACHE İŞLEMLERİ ====================

// GetPolyclinicTypes - Poliklinik tiplerini cache'den getirir
func (cs *CacheService) GetPolyclinicTypes() ([]model.PolyclinicType, error) {
	// Cache'den kontrol et
	cachedData, err := cs.redisClient.Get(cs.ctx, CACHE_POLYCLINIC_TYPES).Result()
	if err == nil {
		// Cache hit - JSON'dan parse et
		var polyclinicTypes []model.PolyclinicType
		if err := json.Unmarshal([]byte(cachedData), &polyclinicTypes); err == nil {
			return polyclinicTypes, nil
		}
	}

	// Cache miss - database'den yükle
	var polyclinicTypes []model.PolyclinicType
	result := database.DB.Order("name ASC").Find(&polyclinicTypes)
	if result.Error != nil {
		return nil, fmt.Errorf("poliklinik tipleri yüklenemedi: %v", result.Error)
	}

	// Cache'e kaydet
	cs.cachePolyclinicTypes(polyclinicTypes)

	return polyclinicTypes, nil
}

// ==================== PRIVATE CACHE HELPER'LARI ====================

// cacheProvinces - İlleri cache'e kaydeder
func (cs *CacheService) cacheProvinces(provinces []model.Province) {
	if jsonData, err := json.Marshal(provinces); err == nil {
		cs.redisClient.Set(cs.ctx, CACHE_PROVINCES, jsonData, cs.defaultTTL)
	}
}

// cacheDistricts - İlçeleri cache'e kaydeder
func (cs *CacheService) cacheDistricts(provinceID uint, districts []model.District) {
	cacheKey := fmt.Sprintf("%s%d", CACHE_DISTRICTS_PREFIX, provinceID)
	if jsonData, err := json.Marshal(districts); err == nil {
		cs.redisClient.Set(cs.ctx, cacheKey, jsonData, cs.defaultTTL)
	}
}

// cacheJobGroups - Meslek gruplarını cache'e kaydeder
func (cs *CacheService) cacheJobGroups(jobGroups []model.JobGroup) {
	if jsonData, err := json.Marshal(jobGroups); err == nil {
		cs.redisClient.Set(cs.ctx, CACHE_JOB_GROUPS, jsonData, cs.defaultTTL)
	}
}

// cacheJobTitles - Unvanları cache'e kaydeder
func (cs *CacheService) cacheJobTitles(jobGroupID uint, jobTitles []model.JobTitle) {
	cacheKey := fmt.Sprintf("%s%d", CACHE_JOB_TITLES, jobGroupID)
	if jsonData, err := json.Marshal(jobTitles); err == nil {
		cs.redisClient.Set(cs.ctx, cacheKey, jsonData, cs.defaultTTL)
	}
}

// cachePolyclinicTypes - Poliklinik tiplerini cache'e kaydeder
func (cs *CacheService) cachePolyclinicTypes(polyclinicTypes []model.PolyclinicType) {
	if jsonData, err := json.Marshal(polyclinicTypes); err == nil {
		cs.redisClient.Set(cs.ctx, CACHE_POLYCLINIC_TYPES, jsonData, cs.defaultTTL)
	}
}

// ==================== CACHE YÖNETİM FONKSİYONLARI ====================

// InvalidateAllMasterData - Tüm master data cache'ini temizler
// Admin panelinde veri güncellendiğinde kullanılır
func (cs *CacheService) InvalidateAllMasterData() error {
	keys := []string{
		CACHE_PROVINCES,
		CACHE_JOB_GROUPS,
		CACHE_POLYCLINIC_TYPES,
	}

	// Pattern'li key'leri de temizle (districts, job_titles)
	districtKeys, _ := cs.redisClient.Keys(cs.ctx, CACHE_DISTRICTS_PREFIX+"*").Result()
	jobTitleKeys, _ := cs.redisClient.Keys(cs.ctx, CACHE_JOB_TITLES+"*").Result()

	keys = append(keys, districtKeys...)
	keys = append(keys, jobTitleKeys...)

	if len(keys) > 0 {
		return cs.redisClient.Del(cs.ctx, keys...).Err()
	}

	return nil
}

// GetCacheStats - Cache istatistiklerini döndürür (monitoring için)
func (cs *CacheService) GetCacheStats() map[string]interface{} {
	stats := make(map[string]interface{})

	// Her cache key'i için TTL bilgisi al
	cacheKeys := []string{
		CACHE_PROVINCES,
		CACHE_JOB_GROUPS,
		CACHE_POLYCLINIC_TYPES,
	}

	for _, key := range cacheKeys {
		ttl := cs.redisClient.TTL(cs.ctx, key).Val()
		exists := cs.redisClient.Exists(cs.ctx, key).Val()
		stats[key] = map[string]interface{}{
			"exists": exists == 1,
			"ttl":    ttl.Seconds(),
		}
	}

	return stats
}
