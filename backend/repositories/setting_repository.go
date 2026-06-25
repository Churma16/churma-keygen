package repositories

import (
	"churma-keygen/backend/models"

	"gorm.io/gorm"
)

type SettingRepository interface {
	Get(key string) (*models.Setting, error)
	Set(key, value string) error
}

type GormSettingRepository struct {
	db *gorm.DB
}

func NewSettingRepository(db *gorm.DB) SettingRepository {
	return &GormSettingRepository{db: db}
}

func (r *GormSettingRepository) Get(key string) (*models.Setting, error) {
	var setting models.Setting
	err := r.db.Where("key = ?", key).First(&setting).Error
	if err != nil {
		return nil, err
	}
	return &setting, nil
}

func (r *GormSettingRepository) Set(key, value string) error {
	var setting models.Setting
	err := r.db.Where("key = ?", key).First(&setting).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			newSetting := models.Setting{
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
