package repository

import (
	"hospital-platform/database"
	"hospital-platform/model"
)

// PolyclinicRepository poliklinik veritabanı işlemlerini yönetir
type PolyclinicRepository struct{}

// NewPolyclinicRepository yeni bir poliklinik repository'si oluşturur
func NewPolyclinicRepository() *PolyclinicRepository {
	return &PolyclinicRepository{}
}

// ==================== POLYCLİNİC TYPES (Master Data) ====================

// GetAllPolyclinicTypes tüm poliklinik türlerini getirir
func (r *PolyclinicRepository) GetAllPolyclinicTypes() ([]model.PolyclinicType, error) {
	var types []model.PolyclinicType
	result := database.DB.Order("name ASC").Find(&types)
	return types, result.Error
}

// GetPolyclinicTypeByID ID'ye göre poliklinik türünü getirir
func (r *PolyclinicRepository) GetPolyclinicTypeByID(id uint) (*model.PolyclinicType, error) {
	var polyclinicType model.PolyclinicType
	result := database.DB.First(&polyclinicType, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &polyclinicType, nil
}

// ==================== HOSPITAL POLYCLİNİCS ====================

// CreateHospitalPolyclinic hastane poliklinik oluşturur
func (r *PolyclinicRepository) CreateHospitalPolyclinic(hospitalPolyclinic *model.HospitalPolyclinic) error {
	result := database.DB.Create(hospitalPolyclinic)
	return result.Error
}

// GetHospitalPolyclinicByID ID'ye göre hastane poliklinik getirir
func (r *PolyclinicRepository) GetHospitalPolyclinicByID(id uint) (*model.HospitalPolyclinic, error) {
	var hospitalPolyclinic model.HospitalPolyclinic
	result := database.DB.Preload("Hospital").Preload("PolyclinicType").First(&hospitalPolyclinic, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &hospitalPolyclinic, nil
}

// CheckHospitalPolyclinicExists hastanede aynı poliklinik türü var mı kontrol eder
func (r *PolyclinicRepository) CheckHospitalPolyclinicExists(hospitalID, polyclinicTypeID uint) (bool, error) {
	var count int64
	result := database.DB.Model(&model.HospitalPolyclinic{}).
		Where("hospital_id = ? AND polyclinic_type_id = ? AND is_active = ?", hospitalID, polyclinicTypeID, true).
		Count(&count)

	if result.Error != nil {
		return false, result.Error
	}

	return count > 0, nil
}

// GetHospitalPolyclinicsSummary hastane polikliniklerini personel sayılarıyla getirir
func (r *PolyclinicRepository) GetHospitalPolyclinicsSummary(hospitalID uint) ([]model.HospitalPolyclinicSummary, error) {
	var summaries []model.HospitalPolyclinicSummary

	// SQL sorgusu: JOIN ile poliklinik bilgileri + personel sayıları
	query := `
		SELECT 
			hp.id,
			pt.name as polyclinic_type_name,
			hp.floor,
			hp.room_number,
			hp.is_active,
			COALESCE(staff_counts.total_staff, 0) as total_staff_count
		FROM hospital_polyclinics hp
		LEFT JOIN polyclinic_types pt ON hp.polyclinic_type_id = pt.id
		LEFT JOIN (
			SELECT 
				polyclinic_id,
				COUNT(*) as total_staff
			FROM staffs 
			WHERE is_active = true 
			GROUP BY polyclinic_id
		) staff_counts ON hp.id = staff_counts.polyclinic_id
		WHERE hp.hospital_id = ? AND hp.is_active = true
		ORDER BY pt.name ASC
	`

	type queryResult struct {
		ID                 uint   `db:"id"`
		PolyclinicTypeName string `db:"polyclinic_type_name"`
		Floor              int    `db:"floor"`
		RoomNumber         int    `db:"room_number"`
		IsActive           bool   `db:"is_active"`
		TotalStaffCount    int    `db:"total_staff_count"`
	}

	var results []queryResult
	err := database.DB.Raw(query, hospitalID).Scan(&results).Error
	if err != nil {
		return nil, err
	}

	// Her poliklinik için meslek grubuna göre personel sayıları
	for _, result := range results {
		summary := model.HospitalPolyclinicSummary{
			ID:                 result.ID,
			PolyclinicTypeName: result.PolyclinicTypeName,
			Floor:              result.Floor,
			RoomNumber:         result.RoomNumber,
			IsActive:           result.IsActive,
			TotalStaffCount:    result.TotalStaffCount,
		}

		// Meslek grubuna göre personel sayıları
		staffByGroup, err := r.getStaffCountByJobGroup(result.ID)
		if err != nil {
			return nil, err
		}
		summary.StaffByJobGroup = staffByGroup

		summaries = append(summaries, summary)
	}

	return summaries, nil
}

// getStaffCountByJobGroup poliklinik için meslek grubuna göre personel sayılarını getirir
func (r *PolyclinicRepository) getStaffCountByJobGroup(polyclinicID uint) ([]model.StaffCountByGroup, error) {
	var counts []model.StaffCountByGroup

	query := `
		SELECT 
			jg.name as job_group_name,
			COUNT(*) as count
		FROM staffs s
		LEFT JOIN job_groups jg ON s.job_group_id = jg.id
		WHERE s.polyclinic_id = ? AND s.is_active = true
		GROUP BY jg.id, jg.name
		ORDER BY jg.name ASC
	`

	err := database.DB.Raw(query, polyclinicID).Scan(&counts).Error
	return counts, err
}

// UpdateHospitalPolyclinic hastane poliklinik günceller
func (r *PolyclinicRepository) UpdateHospitalPolyclinic(hospitalPolyclinic *model.HospitalPolyclinic) error {
	result := database.DB.Save(hospitalPolyclinic)
	return result.Error
}

// DeleteHospitalPolyclinic hastane poliklinik siler (soft delete)
func (r *PolyclinicRepository) DeleteHospitalPolyclinic(id uint) error {
	result := database.DB.Delete(&model.HospitalPolyclinic{}, id)
	return result.Error
}

// ==================== LEGACY POLYCLİNİC (Geriye uyumluluk) ====================

// CreatePolyclinic eski poliklinik oluşturur (legacy)
func CreatePolyclinic(p *model.Polyclinic) error {
	result := database.DB.Create(p)
	return result.Error
}

// GetAllPolyclinics tüm eski poliklinikleri getirir (legacy)
func GetAllPolyclinics() ([]model.Polyclinic, error) {
	var polyclinics []model.Polyclinic
	result := database.DB.Find(&polyclinics)
	return polyclinics, result.Error
}

// GetPolyclinicByID ID'ye göre eski poliklinik getirir (legacy)
func GetPolyclinicByID(id uint) (*model.Polyclinic, error) {
	var poly model.Polyclinic
	result := database.DB.First(&poly, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &poly, nil
}

// UpdatePolyclinic eski poliklinik günceller (legacy)
func UpdatePolyclinic(poly *model.Polyclinic) error {
	result := database.DB.Save(poly)
	return result.Error
}

// DeletePolyclinic eski poliklinik siler (legacy)
func DeletePolyclinic(id uint) error {
	result := database.DB.Delete(&model.Polyclinic{}, id)
	return result.Error
}
