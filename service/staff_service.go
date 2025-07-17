package service

import (
	"encoding/json"
	"fmt"
	"hospital-platform/model"
	"hospital-platform/repository"
)

// StaffService - Personel işlemlerinin iş mantığını koordine eden servis katmanı
// Validasyon kuralları, iş kuralları ve repository'ler arası ilişkileri yönetir
type StaffService struct {
	staffRepo      *repository.StaffRepository      // Personel veritabanı işlemleri
	polyclinicRepo *repository.PolyclinicRepository // Poliklinik doğrulama işlemleri
}

// NewStaffService - Bağımlılıkları enjekte ederek yeni servis instance'ı oluşturur
// Repository'leri initialize eder ve hazır hale getirir
func NewStaffService() *StaffService {
	return &StaffService{
		staffRepo:      repository.NewStaffRepository(),
		polyclinicRepo: repository.NewPolyclinicRepository(),
	}
}

// ==================== PERSONEL YÖNETİM İŞLEMLERİ ====================

// CreateStaff - Yeni personel oluşturma işlemini baştan sona yönetir
// Validasyon, iş kuralları kontrolü ve veritabanına kaydetme işlemlerini yapar
func (s *StaffService) CreateStaff(req *model.CreateStaffRequest, hospitalID uint) (*model.Staff, []model.ValidationError, error) {
	// İlk olarak tüm iş kurallarını ve validasyonları kontrol et
	validationErrors := s.validateCreateStaff(req, hospitalID)
	if len(validationErrors) > 0 {
		return nil, validationErrors, nil
	}

	// 2. WorkDays'i JSON string'e çevir
	workDaysJSON, err := json.Marshal(req.WorkDays)
	if err != nil {
		return nil, nil, fmt.Errorf("çalışma günleri işlenemedi: %v", err)
	}

	// 3. Staff model oluştur
	staff := &model.Staff{
		HospitalID:   hospitalID,
		PolyclinicID: req.PolyclinicID,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		TCKN:         req.TCKN,
		Phone:        req.Phone,
		JobGroupID:   req.JobGroupID,
		JobTitleID:   req.JobTitleID,
		WorkDays:     string(workDaysJSON),
		IsActive:     true,
	}

	// 4. Veritabanına kaydet
	err = s.staffRepo.Create(staff)
	if err != nil {
		return nil, nil, fmt.Errorf("personel oluşturulamadı: %v", err)
	}

	// 5. İlişkilerle beraber geri döndür
	result, err := s.staffRepo.GetByID(staff.ID)
	if err != nil {
		return nil, nil, fmt.Errorf("oluşturulan personel getirilemedi: %v", err)
	}
	return result, nil, nil
}

// GetStaffByID ID'ye göre personel getirir
func (s *StaffService) GetStaffByID(id uint, hospitalID uint) (*model.Staff, error) {
	staff, err := s.staffRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("personel bulunamadı")
	}

	// Hastane sahipliği kontrolü
	if staff.HospitalID != hospitalID {
		return nil, fmt.Errorf("bu personel size ait değil")
	}

	return staff, nil
}

// UpdateStaff personel bilgilerini günceller
func (s *StaffService) UpdateStaff(id uint, req *model.UpdateStaffRequest, hospitalID uint) (*model.Staff, []model.ValidationError, error) {
	// 1. Mevcut personeli getir
	staff, err := s.GetStaffByID(id, hospitalID)
	if err != nil {
		return nil, nil, err
	}

	// 2. Validasyon
	validationErrors := s.validateUpdateStaff(req, hospitalID, &id)
	if len(validationErrors) > 0 {
		return nil, validationErrors, nil
	}

	// 3. WorkDays'i JSON string'e çevir
	workDaysJSON, err := json.Marshal(req.WorkDays)
	if err != nil {
		return nil, nil, fmt.Errorf("çalışma günleri işlenemedi: %v", err)
	}

	// 4. Güncelle
	staff.FirstName = req.FirstName
	staff.LastName = req.LastName
	staff.Phone = req.Phone
	staff.JobGroupID = req.JobGroupID
	staff.JobTitleID = req.JobTitleID
	staff.PolyclinicID = req.PolyclinicID
	staff.WorkDays = string(workDaysJSON)
	staff.IsActive = req.IsActive

	err = s.staffRepo.Update(staff)
	if err != nil {
		return nil, nil, fmt.Errorf("personel güncellenemedi: %v", err)
	}

	return staff, nil, nil
}

// DeleteStaff personeli siler
func (s *StaffService) DeleteStaff(id uint, hospitalID uint) error {
	// 1. Personel hastaneye ait mi kontrol et
	_, err := s.GetStaffByID(id, hospitalID)
	if err != nil {
		return err
	}

	// 2. Sil
	return s.staffRepo.Delete(id)
}

// ==================== LİSTELEME VE FİLTRELEME ====================

// GetStaffList sayfalandırılmış personel listesi getirir
func (s *StaffService) GetStaffList(req *model.StaffListRequest, hospitalID uint) (*model.StaffListResponse, error) {
	// Default değerler
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 || req.PageSize > 100 {
		req.PageSize = 10
	}

	return s.staffRepo.GetPaginatedStaff(hospitalID, req)
}

// ==================== MASTER DATA ====================

// GetJobGroups meslek gruplarını getirir
func (s *StaffService) GetJobGroups() ([]model.JobGroup, error) {
	return s.staffRepo.GetJobGroups()
}

// GetJobTitlesByGroup meslek grubuna göre unvanları getirir
func (s *StaffService) GetJobTitlesByGroup(jobGroupID uint) ([]model.JobTitle, error) {
	return s.staffRepo.GetJobTitlesByGroup(jobGroupID)
}

// ==================== VALİDASYON METHODS ====================

// validateCreateStaff personel oluşturma validasyonu
func (s *StaffService) validateCreateStaff(req *model.CreateStaffRequest, hospitalID uint) []model.ValidationError {
	var errors []model.ValidationError

	// TC kimlik numarası benzersizlik kontrolü
	exists, err := s.staffRepo.CheckTCKNExists(req.TCKN, nil)
	if err == nil && exists {
		errors = append(errors, model.ValidationError{
			Field:   "tc",
			Message: "Bu TC kimlik numarası zaten kullanılıyor",
		})
	}

	// Telefon benzersizlik kontrolü
	exists, err = s.staffRepo.CheckPhoneExists(req.Phone, nil)
	if err == nil && exists {
		errors = append(errors, model.ValidationError{
			Field:   "phone",
			Message: "Bu telefon numarası zaten kullanılıyor",
		})
	}

	// Unvan benzersizlik kontrolü (Başhekim vb.)
	canAssign, err := s.staffRepo.CheckUniqueJobTitle(hospitalID, req.JobTitleID, nil)
	if err == nil && !canAssign {
		errors = append(errors, model.ValidationError{
			Field:   "job_title_id",
			Message: "Bu unvandan hastanede sadece bir tane olabilir",
		})
	}

	// Poliklinik kontrolü
	if req.PolyclinicID != nil {
		polyclinic, err := s.polyclinicRepo.GetHospitalPolyclinicByID(*req.PolyclinicID)
		if err != nil || polyclinic.HospitalID != hospitalID {
			errors = append(errors, model.ValidationError{
				Field:   "polyclinic_id",
				Message: "Geçersiz poliklinik seçimi",
			})
		}
	}

	// Çalışma günleri kontrolü
	if len(req.WorkDays) == 0 {
		errors = append(errors, model.ValidationError{
			Field:   "work_days",
			Message: "En az bir çalışma günü seçilmelidir",
		})
	} else {
		for _, day := range req.WorkDays {
			if day < 1 || day > 7 {
				errors = append(errors, model.ValidationError{
					Field:   "work_days",
					Message: "Geçersiz gün değeri (1-7 arasında olmalı)",
				})
				break
			}
		}
	}

	return errors
}

// validateUpdateStaff personel güncelleme validasyonu
func (s *StaffService) validateUpdateStaff(req *model.UpdateStaffRequest, hospitalID uint, excludeID *uint) []model.ValidationError {
	var errors []model.ValidationError

	// Telefon benzersizlik kontrolü (kendi ID'si hariç)
	exists, err := s.staffRepo.CheckPhoneExists(req.Phone, excludeID)
	if err == nil && exists {
		errors = append(errors, model.ValidationError{
			Field:   "phone",
			Message: "Bu telefon numarası zaten kullanılıyor",
		})
	}

	// Unvan benzersizlik kontrolü (kendi ID'si hariç)
	canAssign, err := s.staffRepo.CheckUniqueJobTitle(hospitalID, req.JobTitleID, excludeID)
	if err == nil && !canAssign {
		errors = append(errors, model.ValidationError{
			Field:   "job_title_id",
			Message: "Bu unvandan hastanede sadece bir tane olabilir",
		})
	}

	// Poliklinik kontrolü
	if req.PolyclinicID != nil {
		polyclinic, err := s.polyclinicRepo.GetHospitalPolyclinicByID(*req.PolyclinicID)
		if err != nil || polyclinic.HospitalID != hospitalID {
			errors = append(errors, model.ValidationError{
				Field:   "polyclinic_id",
				Message: "Geçersiz poliklinik seçimi",
			})
		}
	}

	// Çalışma günleri kontrolü
	if len(req.WorkDays) == 0 {
		errors = append(errors, model.ValidationError{
			Field:   "work_days",
			Message: "En az bir çalışma günü seçilmelidir",
		})
	} else {
		for _, day := range req.WorkDays {
			if day < 1 || day > 7 {
				errors = append(errors, model.ValidationError{
					Field:   "work_days",
					Message: "Geçersiz gün değeri (1-7 arasında olmalı)",
				})
				break
			}
		}
	}

	return errors
}
