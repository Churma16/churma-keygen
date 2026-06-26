package repositories

import (
	"churma-keygen/backend/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GormClientRepository struct {
	db *gorm.DB
}

func NewClientRepository(db *gorm.DB) domain.ClientRepository {
	return &GormClientRepository{db: db}
}

func (r *GormClientRepository) FindAll() ([]domain.Client, error) {
	var clients []domain.Client
	err := r.db.Preload("Licenses").Order("name ASC").Find(&clients).Error
	return clients, err
}

func (r *GormClientRepository) FindByID(id uuid.UUID) (*domain.Client, error) {
	var client domain.Client
	err := r.db.First(&client, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &client, nil
}

func (r *GormClientRepository) Create(client *domain.Client) error {
	return r.db.Create(client).Error
}

func (r *GormClientRepository) Update(client *domain.Client) error {
	return r.db.Save(client).Error
}

func (r *GormClientRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&domain.Client{}, "id = ?", id).Error
}

func (r *GormClientRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&domain.Client{}).Count(&count).Error
	return count, err
}
