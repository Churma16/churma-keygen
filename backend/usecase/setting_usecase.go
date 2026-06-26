package usecase

import (
	"churma-keygen/backend/domain"

	"gorm.io/gorm"
)

type SettingUsecase interface {
	GetSetting(key string) (*domain.Setting, error)
	UpdateSetting(key, value string) (*domain.Setting, error)
}

type settingUsecaseImpl struct {
	settingRepo domain.SettingRepository
}

func NewSettingUsecase(settingRepo domain.SettingRepository) SettingUsecase {
	return &settingUsecaseImpl{settingRepo: settingRepo}
}

func (s *settingUsecaseImpl) GetSetting(key string) (*domain.Setting, error) {
	setting, err := s.settingRepo.Get(key)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &domain.Setting{
				Key:   key,
				Value: "",
			}, nil
		}
		return nil, err
	}
	return setting, nil
}

func (s *settingUsecaseImpl) UpdateSetting(key, value string) (*domain.Setting, error) {
	err := s.settingRepo.Set(key, value)
	if err != nil {
		return nil, err
	}
	return s.settingRepo.Get(key)
}
