package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID           uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	Username     string         `gorm:"type:varchar(50);unique;not null" json:"username"`
	PasswordHash string         `gorm:"type:varchar(255);not null" json:"-"`
	Role         string         `gorm:"type:varchar(20);default:'SUPERADMIN';not null" json:"role"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return
}

type Client struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	Name      string         `gorm:"type:varchar(100);not null" json:"name"`
	OwnerName string         `gorm:"type:varchar(100)" json:"owner_name"`
	Phone     string         `gorm:"type:varchar(20)" json:"phone"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Licenses  []License      `gorm:"constraint:OnDelete:CASCADE;" json:"licenses,omitempty"`
}

func (c *Client) BeforeCreate(tx *gorm.DB) (err error) {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return
}

type License struct {
	ID             uuid.UUID       `gorm:"type:uuid;primaryKey" json:"id"`
	ClientID       uuid.UUID       `gorm:"type:uuid;not null;index" json:"client_id"`
	Client         *Client         `gorm:"foreignKey:ClientID" json:"client,omitempty"`
	LicenseCode    string          `gorm:"type:varchar(50);unique;not null;index" json:"license_code"`
	HardwareID     string          `gorm:"type:varchar(64)" json:"hardware_id"`
	TrialLimit     int             `gorm:"default:100;not null" json:"trial_limit"`
	Status         string          `gorm:"type:varchar(20);default:'UNASSIGNED';not null" json:"status"` // UNASSIGNED, ACTIVE, SUSPENDED, REVOKED
	ExpiresAt      *time.Time      `json:"expires_at"`
	ActivatedAt    *time.Time      `json:"activated_at"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
	DeletedAt      gorm.DeletedAt  `gorm:"index" json:"-"`
	ActivationLogs []ActivationLog `gorm:"constraint:OnDelete:SET NULL;" json:"activation_logs,omitempty"`
}

func (l *License) BeforeCreate(tx *gorm.DB) (err error) {
	if l.ID == uuid.Nil {
		l.ID = uuid.New()
	}
	return
}

type ActivationLog struct {
	ID                uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	LicenseID         *uuid.UUID     `gorm:"type:uuid;index" json:"license_id"`
	License           *License       `gorm:"foreignKey:LicenseID" json:"license,omitempty"`
	AttemptedCode     string         `gorm:"type:varchar(50);not null" json:"attempted_code"`
	HardwareIDAttempt string         `gorm:"type:varchar(64);not null" json:"hardware_id_attempt"`
	IPAddress         string         `gorm:"type:varchar(45);not null" json:"ip_address"`
	Status            string         `gorm:"type:varchar(30);not null" json:"status"` // SUCCESS, INVALID_KEY, HWID_MISMATCH, SUSPENDED_KEY
	CreatedAt         time.Time      `json:"created_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`
}

func (al *ActivationLog) BeforeCreate(tx *gorm.DB) (err error) {
	if al.ID == uuid.Nil {
		al.ID = uuid.New()
	}
	return
}

type Setting struct {
	Key       string    `gorm:"type:varchar(50);primaryKey" json:"key"`
	Value     string    `gorm:"type:text;not null" json:"value"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

