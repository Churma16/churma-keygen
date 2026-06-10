package repositories

import (
	"churma-keygen/backend/models"

	"gorm.io/gorm"
)

type ActivationLogRepository interface {
	FindAll(limit int) ([]models.ActivationLog, error)
	Create(log *models.ActivationLog) error
}

type GormActivationLogRepository struct {
	db *gorm.DB
}

func NewActivationLogRepository(db *gorm.DB) ActivationLogRepository {
	return &GormActivationLogRepository{db: db}
}

func (r *GormActivationLogRepository) FindAll(limit int) ([]models.ActivationLog, error) {
	var logs []models.ActivationLog
	err := r.db.Preload("License").Preload("License.Client").Order("created_at DESC").Limit(limit).Find(&logs).Error
	return logs, err
}

func (r *GormActivationLogRepository) Create(log *models.ActivationLog) error {
	return r.db.Create(log).Error
}
