package repositories

import (
	"churma-keygen/backend/domain"

	"gorm.io/gorm"
)

type GormSettingRepository struct {
	db *gorm.DB
}

func NewSettingRepository(db *gorm.DB) domain.SettingRepository {
	return &GormSettingRepository{db: db}
}

func (r *GormSettingRepository) Get(key string) (*domain.Setting, error) {
	var setting domain.Setting
	err := r.db.Where("key = ?", key).First(&setting).Error
	if err != nil {
		return nil, err
	}
	return &setting, nil
}

func (r *GormSettingRepository) Set(key, value string) error {
	var setting domain.Setting
	err := r.db.Where("key = ?", key).First(&setting).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			newSetting := domain.Setting{
				Key:   key,
				Value: value,
			}
			return r.db.Create(&newSetting).Error
		}
		return err
	}
	setting.Value = value
	return r.db.Save(&setting).Error
}
