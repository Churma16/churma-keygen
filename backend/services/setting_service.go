package services

import (
	"churma-keygen/backend/models"
	"churma-keygen/backend/repositories"
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
	return s.settingRepo.Get(key)
}

func (s *settingServiceImpl) UpdateSetting(key, value string) (*models.Setting, error) {
	err := s.settingRepo.Set(key, value)
	if err != nil {
		return nil, err
	}
	return s.settingRepo.Get(key)
}
