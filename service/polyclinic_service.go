package service

import (
	"fmt"
	"hospital-platform/model"
	"hospital-platform/repository"
)

// PolyclinicService poliklinik iş mantığını yönetir
type PolyclinicService struct {
	polyclinicRepo *repository.PolyclinicRepository
	locationRepo   *repository.LocationRepository
}

// NewPolyclinicService yeni bir poliklinik servisi oluşturur
func NewPolyclinicService() *PolyclinicService {
	return &PolyclinicService{
		polyclinicRepo: repository.NewPolyclinicRepository(),
		locationRepo:   repository.NewLocationRepository(),
	}
}

// GetAllPolyclinicTypes master data poliklinik türlerini getirir
func (s *PolyclinicService) GetAllPolyclinicTypes() ([]model.PolyclinicType, error) {
	return s.polyclinicRepo.GetAllPolyclinicTypes()
}

// AddPolyclinicToHospital hastaneye yeni poliklinik ekler
func (s *PolyclinicService) AddPolyclinicToHospital(req *model.AddPolyclinicRequest, hospitalID uint) (*model.HospitalPolyclinic, error) {
	// 1. Poliklinik türünün var olup olmadığını kontrol et
	_, err := s.polyclinicRepo.GetPolyclinicTypeByID(req.PolyclinicTypeID)
	if err != nil {
		return nil, fmt.Errorf("geçersiz poliklinik türü")
	}

	// 2. Aynı poliklinik türü bu hastanede zaten var mı kontrol et
	exists, err := s.polyclinicRepo.CheckHospitalPolyclinicExists(hospitalID, req.PolyclinicTypeID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, fmt.Errorf("bu poliklinik türü zaten hastanede mevcut")
	}

	// 3. Yeni hastane poliklinik oluştur
	hospitalPolyclinic := &model.HospitalPolyclinic{
		HospitalID:       hospitalID,
		PolyclinicTypeID: req.PolyclinicTypeID,
		Floor:            req.Floor,
		RoomNumber:       req.RoomNumber,
		IsActive:         true,
	}

	err = s.polyclinicRepo.CreateHospitalPolyclinic(hospitalPolyclinic)
	if err != nil {
		return nil, fmt.Errorf("poliklinik eklenemedi: %v", err)
	}

	// 4. İlişkilerle beraber geri döndür
	return s.polyclinicRepo.GetHospitalPolyclinicByID(hospitalPolyclinic.ID)
}

// GetHospitalPolyclinics hastaneye ait poliklinikleri temel bilgilerle getirir
func (s *PolyclinicService) GetHospitalPolyclinics(hospitalID uint) ([]model.HospitalPolyclinicSummary, error) {
	return s.polyclinicRepo.GetHospitalPolyclinicsSummary(hospitalID)
}

// UpdateHospitalPolyclinic hastane poliklinik bilgilerini günceller
func (s *PolyclinicService) UpdateHospitalPolyclinic(id uint, req *model.UpdatePolyclinicRequest, hospitalID uint) (*model.HospitalPolyclinic, error) {
	// 1. Poliklinik hastaneye ait mi kontrol et
	polyclinic, err := s.polyclinicRepo.GetHospitalPolyclinicByID(id)
	if err != nil {
		return nil, fmt.Errorf("poliklinik bulunamadı")
	}

	if polyclinic.HospitalID != hospitalID {
		return nil, fmt.Errorf("bu poliklinik size ait değil")
	}

	// 2. Güncelle
	polyclinic.Floor = req.Floor
	polyclinic.RoomNumber = req.RoomNumber
	polyclinic.IsActive = req.IsActive

	err = s.polyclinicRepo.UpdateHospitalPolyclinic(polyclinic)
	if err != nil {
		return nil, fmt.Errorf("poliklinik güncellenemedi: %v", err)
	}

	return polyclinic, nil
}

// DeleteHospitalPolyclinic hastane poliklinik siler
func (s *PolyclinicService) DeleteHospitalPolyclinic(id uint, hospitalID uint) error {
	// 1. Poliklinik hastaneye ait mi kontrol et
	polyclinic, err := s.polyclinicRepo.GetHospitalPolyclinicByID(id)
	if err != nil {
		return fmt.Errorf("poliklinik bulunamadı")
	}

	if polyclinic.HospitalID != hospitalID {
		return fmt.Errorf("bu poliklinik size ait değil")
	}

	// 2. Sil
	return s.polyclinicRepo.DeleteHospitalPolyclinic(id)
}
