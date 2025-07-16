package service

import (
	"errors"
	"hospital-platform/model"
	"hospital-platform/repository"
)

// Yeni poliklinik oluştur
func CreatePolyclinicService(p *model.Polyclinic) error {
	if p.Name == "" {
		return errors.New("poliklinik adı boş olamaz")
	}
	return repository.CreatePolyclinic(p)
}

// Tüm poliklinikleri getir
func GetAllPolyclinicsService() ([]model.Polyclinic, error) {
	return repository.GetAllPolyclinics()
}

// ID'ye göre poliklinik getir
func GetPolyclinicByIDService(id uint) (*model.Polyclinic, error) {
	return repository.GetPolyclinicByID(id)
}

// Silme işlemi
func DeletePolyclinicService(id uint) error {
	return repository.DeletePolyclinic(id)
}

func UpdatePolyclinicService(id uint, updated *model.Polyclinic) error {
	existing, err := repository.GetPolyclinicByID(id)
	if err != nil {
		return err
	}

	existing.Name = updated.Name

	return repository.UpdatePolyclinic(existing)
}
