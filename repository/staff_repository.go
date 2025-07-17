package repository

import (
	"encoding/json"
	"fmt"
	"hospital-platform/database"
	"hospital-platform/model"
	"strings"
)

// StaffRepository - Personel verilerine erişim katmanı
// Veritabanı ile personel CRUD işlemlerini gerçekleştiren ana sınıf
type StaffRepository struct{}

// NewStaffRepository - Yeni bir personel repository nesnesi oluşturur
// Her çağrıldığında temiz bir instance döner
func NewStaffRepository() *StaffRepository {
	return &StaffRepository{}
}

// ==================== TEMEL VERİTABANI İŞLEMLERİ ====================

// Create - Yeni bir personel kaydını veritabanına ekler
// Çalışma günlerini JSON formatında saklar ve gerekli validasyonları yapar
func (r *StaffRepository) Create(staff *model.Staff) error {
	// Çalışma günleri dizisini JSON string formatına çevir (örn: [1,2,3,4,5])
	workDaysJSON, err := json.Marshal(staff.WorkDays)
	if err != nil {
		return fmt.Errorf("çalışma günleri JSON formatına çevrilemedi: %v", err)
	}
	staff.WorkDays = string(workDaysJSON)

	result := database.DB.Create(staff)
	return result.Error
}

// GetByID - Verilen ID'ye sahip personeli tüm ilişkili verilerle beraber getirir
// Hastane, poliklinik, meslek grubu ve unvan bilgilerini de yükler
func (r *StaffRepository) GetByID(id uint) (*model.Staff, error) {
	var staff model.Staff
	result := database.DB.Preload("Hospital").Preload("Polyclinic").Preload("JobGroup").Preload("JobTitle").First(&staff, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &staff, nil
}

// Update personel bilgilerini günceller
func (r *StaffRepository) Update(staff *model.Staff) error {
	result := database.DB.Save(staff)
	return result.Error
}

// Delete personeli soft delete yapar
func (r *StaffRepository) Delete(id uint) error {
	result := database.DB.Delete(&model.Staff{}, id)
	return result.Error
}

// ==================== VERİFİCATİON METHODS ====================

// CheckTCKNExists TC kimlik numarası var mı kontrol eder
func (r *StaffRepository) CheckTCKNExists(tckn string, excludeID *uint) (bool, error) {
	var count int64
	query := database.DB.Model(&model.Staff{}).Where("tckn = ?", tckn)

	if excludeID != nil {
		query = query.Where("id != ?", *excludeID)
	}

	result := query.Count(&count)
	return count > 0, result.Error
}

// CheckPhoneExists telefon numarası var mı kontrol eder
func (r *StaffRepository) CheckPhoneExists(phone string, excludeID *uint) (bool, error) {
	var count int64
	query := database.DB.Model(&model.Staff{}).Where("phone = ?", phone)

	if excludeID != nil {
		query = query.Where("id != ?", *excludeID)
	}

	result := query.Count(&count)
	return count > 0, result.Error
}

// CheckUniqueJobTitle unvan benzersizliği kontrol eder (örn: Başhekim)
func (r *StaffRepository) CheckUniqueJobTitle(hospitalID, jobTitleID uint, excludeID *uint) (bool, error) {
	// Önce bu unvan unique mi kontrol et
	var jobTitle model.JobTitle
	result := database.DB.First(&jobTitle, jobTitleID)
	if result.Error != nil {
		return false, result.Error
	}

	// Eğer unique değilse kontrole gerek yok
	if !jobTitle.IsUnique {
		return true, nil
	}

	// Unique ise hastanede başka biri var mı bak
	var count int64
	query := database.DB.Model(&model.Staff{}).
		Where("hospital_id = ? AND job_title_id = ? AND is_active = ?", hospitalID, jobTitleID, true)

	if excludeID != nil {
		query = query.Where("id != ?", *excludeID)
	}

	result = query.Count(&count)
	return count == 0, result.Error // 0 ise ekleyebilir
}

// ==================== LİSTELEME VE FİLTRELEME ====================

// GetPaginatedStaff sayfalandırılmış ve filtreli personel listesi getirir
func (r *StaffRepository) GetPaginatedStaff(hospitalID uint, req *model.StaffListRequest) (*model.StaffListResponse, error) {
	var staffList []model.StaffSummary
	var totalCount int64

	// Base query
	baseQuery := r.buildStaffQuery(hospitalID, req)

	// Toplam kayıt sayısını al
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM (%s) as count_query", baseQuery)
	err := database.DB.Raw(countQuery, r.buildQueryParams(hospitalID, req)...).Scan(&totalCount).Error
	if err != nil {
		return nil, fmt.Errorf("toplam kayıt sayısı hesaplanamadı: %v", err)
	}

	// Sayfalama hesaplamaları
	offset := (req.Page - 1) * req.PageSize
	totalPages := int((totalCount + int64(req.PageSize) - 1) / int64(req.PageSize))

	// Ana sorgu - sayfalama ile
	mainQuery := fmt.Sprintf("%s ORDER BY s.first_name ASC, s.last_name ASC LIMIT %d OFFSET %d",
		baseQuery, req.PageSize, offset)

	err = database.DB.Raw(mainQuery, r.buildQueryParams(hospitalID, req)...).Scan(&staffList).Error
	if err != nil {
		return nil, fmt.Errorf("personel listesi getirilemedi: %v", err)
	}

	// WorkDays metinlerini düzenle
	for i := range staffList {
		staffList[i].WorkDaysText = r.formatWorkDays(staffList[i].WorkDaysText)
	}

	// Pagination bilgileri
	pagination := model.PaginationInfo{
		CurrentPage:  req.Page,
		PageSize:     req.PageSize,
		TotalRecords: totalCount,
		TotalPages:   totalPages,
		HasNext:      req.Page < totalPages,
		HasPrev:      req.Page > 1,
	}

	return &model.StaffListResponse{
		Data:       staffList,
		Pagination: pagination,
	}, nil
}

// buildStaffQuery personel sorgusu oluşturur
func (r *StaffRepository) buildStaffQuery(hospitalID uint, req *model.StaffListRequest) string {
	query := `
		SELECT 
			s.id,
			s.first_name,
			s.last_name,
			s.tckn,
			s.phone,
			jg.name as job_group_name,
			jt.name as job_title_name,
			pt.name as polyclinic_type_name,
			s.work_days as work_days_text,
			s.is_active
		FROM staffs s
		LEFT JOIN job_groups jg ON s.job_group_id = jg.id
		LEFT JOIN job_titles jt ON s.job_title_id = jt.id
		LEFT JOIN hospital_polyclinics hp ON s.polyclinic_id = hp.id
		LEFT JOIN polyclinic_types pt ON hp.polyclinic_type_id = pt.id
		WHERE s.hospital_id = ? AND s.deleted_at IS NULL
	`

	// Filtreleme koşulları ekle
	if req.FirstName != "" {
		query += " AND s.first_name ILIKE ?"
	}
	if req.LastName != "" {
		query += " AND s.last_name ILIKE ?"
	}
	if req.TCKN != "" {
		query += " AND s.tckn ILIKE ?"
	}
	if req.JobGroupID != nil {
		query += " AND s.job_group_id = ?"
	}
	if req.JobTitleID != nil {
		query += " AND s.job_title_id = ?"
	}
	if req.PolyclinicID != nil {
		query += " AND s.polyclinic_id = ?"
	}
	if req.IsActive != nil {
		query += " AND s.is_active = ?"
	}

	return query
}

// buildQueryParams sorgu parametrelerini oluşturur
func (r *StaffRepository) buildQueryParams(hospitalID uint, req *model.StaffListRequest) []interface{} {
	params := []interface{}{hospitalID}

	if req.FirstName != "" {
		params = append(params, "%"+req.FirstName+"%")
	}
	if req.LastName != "" {
		params = append(params, "%"+req.LastName+"%")
	}
	if req.TCKN != "" {
		params = append(params, "%"+req.TCKN+"%")
	}
	if req.JobGroupID != nil {
		params = append(params, *req.JobGroupID)
	}
	if req.JobTitleID != nil {
		params = append(params, *req.JobTitleID)
	}
	if req.PolyclinicID != nil {
		params = append(params, *req.PolyclinicID)
	}
	if req.IsActive != nil {
		params = append(params, *req.IsActive)
	}

	return params
}

// formatWorkDays çalışma günlerini okunabilir metne çevirir
func (r *StaffRepository) formatWorkDays(workDaysJSON string) string {
	var days []int
	if err := json.Unmarshal([]byte(workDaysJSON), &days); err != nil {
		return "Bilinmiyor"
	}

	dayNames := map[int]string{
		1: "Pazartesi", 2: "Salı", 3: "Çarşamba", 4: "Perşembe",
		5: "Cuma", 6: "Cumartesi", 7: "Pazar",
	}

	var dayTexts []string
	for _, day := range days {
		if name, exists := dayNames[day]; exists {
			dayTexts = append(dayTexts, name)
		}
	}

	if len(dayTexts) == 0 {
		return "Çalışma günü yok"
	}

	return strings.Join(dayTexts, ", ")
}

// ==================== MASTER DATA HELPERS ====================

// GetJobGroups meslek gruplarını getirir
func (r *StaffRepository) GetJobGroups() ([]model.JobGroup, error) {
	var jobGroups []model.JobGroup
	result := database.DB.Order("name ASC").Find(&jobGroups)
	return jobGroups, result.Error
}

// GetJobTitlesByGroup meslek grubuna göre unvanları getirir
func (r *StaffRepository) GetJobTitlesByGroup(jobGroupID uint) ([]model.JobTitle, error) {
	var jobTitles []model.JobTitle
	result := database.DB.Where("job_group_id = ?", jobGroupID).Order("name ASC").Find(&jobTitles)
	return jobTitles, result.Error
}
