package domain

import (
	"crypto/rand"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

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

// GenerateRandomCode creates a random license code
func (l *License) GenerateRandomCode() {
	const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, 12)
	bytes := make([]byte, 12)
	_, _ = rand.Read(bytes)
	for i, b := range bytes {
		result[i] = chars[b%byte(len(chars))]
	}
	l.LicenseCode = fmt.Sprintf("SFA-%s-%s-%s", string(result[0:4]), string(result[4:8]), string(result[8:12]))
}

// ValidateAndActivate checks if the license can be activated with the given hardware ID
func (l *License) ValidateAndActivate(hardwareID string) (string, error) {
	// Check if suspended or revoked
	if l.Status == "SUSPENDED" || l.Status == "REVOKED" {
		return "SUSPENDED_KEY", errors.New("this license key is suspended or revoked")
	}

	// Check if expired
	if l.ExpiresAt != nil && l.ExpiresAt.Before(time.Now()) {
		return "SUSPENDED_KEY", errors.New("this license key has expired")
	}

	// Handle Hardware ID binding
	if l.HardwareID == "" {
		// First-time activation: bind the HWID
		l.HardwareID = hardwareID
		l.Status = "ACTIVE"
		now := time.Now()
		l.ActivatedAt = &now
		return "SUCCESS", nil
	} else {
		// Subsequent activation: verify HWID
		if l.HardwareID != hardwareID {
			return "HWID_MISMATCH", errors.New("hardware ID mismatch. This license is bound to another machine")
		}
		return "SUCCESS", nil
	}
}

// UpdateStatus updates the license status safely
func (l *License) UpdateStatus(newStatus string) error {
	validStatuses := map[string]bool{
		"UNASSIGNED": true,
		"ACTIVE":     true,
		"SUSPENDED":  true,
		"REVOKED":    true,
	}
	if !validStatuses[newStatus] {
		return errors.New("invalid status. Must be UNASSIGNED, ACTIVE, SUSPENDED, or REVOKED")
	}

	l.Status = newStatus
	if newStatus == "UNASSIGNED" {
		l.HardwareID = ""
		l.ActivatedAt = nil
	}
	return nil
}
