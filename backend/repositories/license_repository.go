package repositories

import (
	"churma-keygen/backend/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GormLicenseRepository struct {
	db *gorm.DB
}

func NewLicenseRepository(db *gorm.DB) domain.LicenseRepository {
	return &GormLicenseRepository{db: db}
}

func (r *GormLicenseRepository) FindAll() ([]domain.License, error) {
	var licenses []domain.License
	err := r.db.Preload("Client").Order("created_at DESC").Find(&licenses).Error
	return licenses, err
}

func (r *GormLicenseRepository) FindByID(id uuid.UUID) (*domain.License, error) {
	var license domain.License
	err := r.db.Preload("Client").First(&license, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &license, nil
}

func (r *GormLicenseRepository) FindByCode(code string) (*domain.License, error) {
	var license domain.License
	err := r.db.Preload("Client").Where("license_code = ?", code).First(&license).Error
	if err != nil {
		return nil, err
	}
	return &license, nil
}

func (r *GormLicenseRepository) Create(license *domain.License) error {
	return r.db.Create(license).Error
}

func (r *GormLicenseRepository) Update(license *domain.License) error {
	return r.db.Save(license).Error
}

func (r *GormLicenseRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&domain.License{}, "id = ?", id).Error
}

func (r *GormLicenseRepository) CountByStatus(status string) (int64, error) {
	var count int64
	err := r.db.Model(&domain.License{}).Where("status = ?", status).Count(&count).Error
	return count, err
}
