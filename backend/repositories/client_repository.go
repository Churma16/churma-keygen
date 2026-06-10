package repositories

import (
	"churma-keygen/backend/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ClientRepository interface {
	FindAll() ([]models.Client, error)
	FindByID(id uuid.UUID) (*models.Client, error)
	Create(client *models.Client) error
	Update(client *models.Client) error
	Delete(id uuid.UUID) error
	Count() (int64, error)
}

type GormClientRepository struct {
	db *gorm.DB
}

func NewClientRepository(db *gorm.DB) ClientRepository {
	return &GormClientRepository{db: db}
}

func (r *GormClientRepository) FindAll() ([]models.Client, error) {
	var clients []models.Client
	err := r.db.Preload("Licenses").Order("name ASC").Find(&clients).Error
	return clients, err
}

func (r *GormClientRepository) FindByID(id uuid.UUID) (*models.Client, error) {
	var client models.Client
	err := r.db.First(&client, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &client, nil
}

func (r *GormClientRepository) Create(client *models.Client) error {
	return r.db.Create(client).Error
}

func (r *GormClientRepository) Update(client *models.Client) error {
	return r.db.Save(client).Error
}

func (r *GormClientRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Client{}, "id = ?", id).Error
}

func (r *GormClientRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&models.Client{}).Count(&count).Error
	return count, err
}
