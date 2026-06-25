package services

import (
	"churma-keygen/backend/models"
	"churma-keygen/backend/repositories"

	"gorm.io/gorm"
)

type SettingService interface {
	GetSetting(key string) (*models.Setting, error)
	UpdateSetting(key, value string) (*models.Setting, error)
}

type settingServiceImpl struct {
	settingRepo repositories.SettingRepository
}

func NewSettingService(settingRepo repositories.SettingRepository) SettingService {
	return &settingServiceImpl{settingRepo: settingRepo}
}

func (s *settingServiceImpl) GetSetting(key string) (*models.Setting, error) {
	setting, err := s.settingRepo.Get(key)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &models.Setting{
				Key:   key,
				Value: "",
			}, nil
		}
		return nil, err
	}
	return setting, nil
}

func (s *settingServiceImpl) UpdateSetting(key, value string) (*models.Setting, error) {
	err := s.settingRepo.Set(key, value)
	if err != nil {
		return nil, err
	}
	return s.settingRepo.Get(key)
}
