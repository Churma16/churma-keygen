package domain

import (
	"github.com/google/uuid"
)

type ClientRepository interface {
	FindAll() ([]Client, error)
	FindByID(id uuid.UUID) (*Client, error)
	Create(client *Client) error
	Update(client *Client) error
	Delete(id uuid.UUID) error
	Count() (int64, error)
}

type LicenseRepository interface {
	FindAll() ([]License, error)
	FindByCode(code string) (*License, error)
	FindByID(id uuid.UUID) (*License, error)
	Create(license *License) error
	Update(license *License) error
	Delete(id uuid.UUID) error
	CountByStatus(status string) (int64, error)
}

type ActivationLogRepository interface {
	FindAll(limit int) ([]ActivationLog, error)
	Create(log *ActivationLog) error
}

type SettingRepository interface {
	Get(key string) (*Setting, error)
	Set(key, value string) error
}

