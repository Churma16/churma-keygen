package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

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

// NewActivationLog creates a new activation log entry
func NewActivationLog(licenseID *uuid.UUID, attemptedCode, hardwareID, ipAddress, status string) *ActivationLog {
	return &ActivationLog{
		ID:                uuid.New(),
		LicenseID:         licenseID,
		AttemptedCode:     attemptedCode,
		HardwareIDAttempt: hardwareID,
		IPAddress:         ipAddress,
		Status:            status,
		CreatedAt:         time.Now(),
	}
}
