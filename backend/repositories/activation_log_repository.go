package repositories

import (
	"churma-keygen/backend/domain"

	"gorm.io/gorm"
)

type GormActivationLogRepository struct {
	db *gorm.DB
}

func NewActivationLogRepository(db *gorm.DB) domain.ActivationLogRepository {
	return &GormActivationLogRepository{db: db}
}

func (r *GormActivationLogRepository) FindAll(limit int) ([]domain.ActivationLog, error) {
	var logs []domain.ActivationLog
	err := r.db.Preload("License").Preload("License.Client").Order("created_at DESC").Limit(limit).Find(&logs).Error
	return logs, err
}

func (r *GormActivationLogRepository) Create(log *domain.ActivationLog) error {
	return r.db.Create(log).Error
}
