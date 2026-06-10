package repositories

import (
	"churma-keygen/backend/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LicenseRepository interface {
	FindAll() ([]models.License, error)
	FindByID(id uuid.UUID) (*models.License, error)
	FindByCode(code string) (*models.License, error)
	Create(license *models.License) error
	Update(license *models.License) error
	Delete(id uuid.UUID) error
	CountByStatus(status string) (int64, error)
}

type GormLicenseRepository struct {
	db *gorm.DB
}

func NewLicenseRepository(db *gorm.DB) LicenseRepository {
	return &GormLicenseRepository{db: db}
}

func (r *GormLicenseRepository) FindAll() ([]models.License, error) {
	var licenses []models.License
	err := r.db.Preload("Client").Order("created_at DESC").Find(&licenses).Error
	return licenses, err
}

func (r *GormLicenseRepository) FindByID(id uuid.UUID) (*models.License, error) {
	var license models.License
	err := r.db.Preload("Client").First(&license, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &license, nil
}

func (r *GormLicenseRepository) FindByCode(code string) (*models.License, error) {
	var license models.License
	err := r.db.Preload("Client").Where("license_code = ?", code).First(&license).Error
	if err != nil {
		return nil, err
	}
	return &license, nil
}

func (r *GormLicenseRepository) Create(license *models.License) error {
	return r.db.Create(license).Error
}

func (r *GormLicenseRepository) Update(license *models.License) error {
	return r.db.Save(license).Error
}

func (r *GormLicenseRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.License{}, "id = ?", id).Error
}

func (r *GormLicenseRepository) CountByStatus(status string) (int64, error) {
	var count int64
	err := r.db.Model(&models.License{}).Where("status = ?", status).Count(&count).Error
	return count, err
}
