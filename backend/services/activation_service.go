package services

import (
	"errors"
	"time"

	"churma-keygen/backend/crypto"
	"churma-keygen/backend/dtos"
	"churma-keygen/backend/models"
	"churma-keygen/backend/repositories"

	"github.com/google/uuid"
)

type ActivationService interface {
	Activate(req dtos.ActivateRequest, ipAddress string) (*dtos.ActivateResponse, error)
	GetPublicKey() (string, error)
	GetContact() *dtos.ContactResponse
}

type activationServiceImpl struct {
	licenseRepo       repositories.LicenseRepository
	activationLogRepo repositories.ActivationLogRepository
	settingRepo       repositories.SettingRepository
}

func NewActivationService(
	licenseRepo repositories.LicenseRepository,
	activationLogRepo repositories.ActivationLogRepository,
	settingRepo repositories.SettingRepository,
) ActivationService {
	return &activationServiceImpl{
		licenseRepo:       licenseRepo,
		activationLogRepo: activationLogRepo,
		settingRepo:       settingRepo,
	}
}

func (s *activationServiceImpl) Activate(req dtos.ActivateRequest, ipAddress string) (*dtos.ActivateResponse, error) {
	license, err := s.licenseRepo.FindByCode(req.LicenseCode)
	if err != nil {
		// Log attempt as INVALID_KEY
		s.logAttempt(nil, req.LicenseCode, req.HardwareID, ipAddress, "INVALID_KEY")
		return nil, errors.New("invalid license code")
	}

	// Check if suspended or revoked
	if license.Status == "SUSPENDED" || license.Status == "REVOKED" {
		s.logAttempt(&license.ID, req.LicenseCode, req.HardwareID, ipAddress, "SUSPENDED_KEY")
		return nil, errors.New("this license key is suspended or revoked")
	}

	// Check if expired
	if license.ExpiresAt != nil && license.ExpiresAt.Before(time.Now()) {
		s.logAttempt(&license.ID, req.LicenseCode, req.HardwareID, ipAddress, "SUSPENDED_KEY")
		return nil, errors.New("this license key has expired")
	}

	// Handle Hardware ID binding
	if license.HardwareID == "" {
		// First-time activation: bind the HWID
		license.HardwareID = req.HardwareID
		license.Status = "ACTIVE"
		now := time.Now()
		license.ActivatedAt = &now

		if err := s.licenseRepo.Update(license); err != nil {
			return nil, errors.New("failed to update license details")
		}
		
		s.logAttempt(&license.ID, req.LicenseCode, req.HardwareID, ipAddress, "SUCCESS")
	} else {
		// Subsequent activation: verify HWID
		if license.HardwareID != req.HardwareID {
			s.logAttempt(&license.ID, req.LicenseCode, req.HardwareID, ipAddress, "HWID_MISMATCH")
			return nil, errors.New("hardware ID mismatch. This license is bound to another machine")
		}

		s.logAttempt(&license.ID, req.LicenseCode, req.HardwareID, ipAddress, "SUCCESS")
	}

	// Generate signed JWT token
	clientName := ""
	if license.Client != nil {
		clientName = license.Client.Name
	}
	token, err := crypto.SignLicenseToken(
		license.ClientID.String(),
		clientName,
		license.HardwareID,
		license.TrialLimit,
		license.ExpiresAt,
	)
	if err != nil {
		return nil, errors.New("failed to sign license token")
	}

	return &dtos.ActivateResponse{
		Token:      token,
		Status:     license.Status,
		ExpiresAt:  license.ExpiresAt,
		TrialLimit: license.TrialLimit,
		ClientName: clientName,
	}, nil
}

func (s *activationServiceImpl) GetPublicKey() (string, error) {
	pubPEM, err := crypto.GetPublicKeyPEM()
	if err != nil {
		return "", errors.New("RSA Public Key is not configured")
	}
	return pubPEM, nil
}

func (s *activationServiceImpl) GetContact() *dtos.ContactResponse {
	var phone string
	setting, err := s.settingRepo.Get("contact_whatsapp")
	if err == nil && setting != nil {
		phone = setting.Value
	}

	var email string
	emailSetting, err := s.settingRepo.Get("contact_email")
	if err == nil && emailSetting != nil {
		email = emailSetting.Value
	}

	var waURL string
	if phone != "" {
		sanitized := sanitizeWhatsAppNumber(phone)
		waURL = "https://wa.me/" + sanitized
	}

	return &dtos.ContactResponse{
		Phone:       phone,
		WhatsAppURL: waURL,
		Email:       email,
	}
}

func sanitizeWhatsAppNumber(phone string) string {
	var digits []rune
	for _, r := range phone {
		if r >= '0' && r <= '9' {
			digits = append(digits, r)
		}
	}
	s := string(digits)
	if len(s) > 0 && s[0] == '0' {
		s = "62" + s[1:]
	}
	return s
}

func (s *activationServiceImpl) logAttempt(licenseID *uuid.UUID, attemptedCode, hardwareID, ipAddress, status string) {
	log := models.ActivationLog{
		ID:                 uuid.New(),
		LicenseID:          licenseID,
		AttemptedCode:      attemptedCode,
		HardwareIDAttempt:  hardwareID,
		IPAddress:          ipAddress,
		Status:             status,
		CreatedAt:          time.Now(),
	}
	_ = s.activationLogRepo.Create(&log)
}
