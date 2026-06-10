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
}

type activationServiceImpl struct {
	licenseRepo       repositories.LicenseRepository
	activationLogRepo repositories.ActivationLogRepository
}

func NewActivationService(
	licenseRepo repositories.LicenseRepository,
	activationLogRepo repositories.ActivationLogRepository,
) ActivationService {
	return &activationServiceImpl{
		licenseRepo:       licenseRepo,
		activationLogRepo: activationLogRepo,
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
