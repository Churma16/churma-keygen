package usecase

import (
	"errors"

	"churma-keygen/backend/crypto"
	"churma-keygen/backend/domain"
	"churma-keygen/backend/dtos"

	"github.com/google/uuid"
)

type ActivationUsecase interface {
	Activate(req dtos.ActivateRequest, ipAddress string) (*dtos.ActivateResponse, error)
	GetPublicKey() (string, error)
	GetContact() *dtos.ContactResponse
}

type activationUsecaseImpl struct {
	licenseRepo       domain.LicenseRepository
	activationLogRepo domain.ActivationLogRepository
	settingRepo       domain.SettingRepository
}

func NewActivationUsecase(
	licenseRepo domain.LicenseRepository,
	activationLogRepo domain.ActivationLogRepository,
	settingRepo domain.SettingRepository,
) ActivationUsecase {
	return &activationUsecaseImpl{
		licenseRepo:       licenseRepo,
		activationLogRepo: activationLogRepo,
		settingRepo:       settingRepo,
	}
}

func (uc *activationUsecaseImpl) Activate(req dtos.ActivateRequest, ipAddress string) (*dtos.ActivateResponse, error) {
	// 1. Fetch from repository
	license, err := uc.licenseRepo.FindByCode(req.LicenseCode)
	if err != nil {
		uc.logAttempt(nil, req.LicenseCode, req.HardwareID, ipAddress, "INVALID_KEY")
		return nil, errors.New("invalid license code")
	}

	// 2. Delegate business logic to Domain
	statusMsg, err := license.ValidateAndActivate(req.HardwareID)
	if err != nil {
		uc.logAttempt(&license.ID, req.LicenseCode, req.HardwareID, ipAddress, statusMsg)
		return nil, err
	}

	// 3. Save domain state back to DB if successful
	if statusMsg == "SUCCESS" {
		if err := uc.licenseRepo.Update(license); err != nil {
			return nil, errors.New("failed to update license details")
		}
	}

	uc.logAttempt(&license.ID, req.LicenseCode, req.HardwareID, ipAddress, "SUCCESS")

	// 4. Generate JWT
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

func (uc *activationUsecaseImpl) GetPublicKey() (string, error) {
	pubPEM, err := crypto.GetPublicKeyPEM()
	if err != nil {
		return "", errors.New("RSA Public Key is not configured")
	}
	return pubPEM, nil
}

func (uc *activationUsecaseImpl) GetContact() *dtos.ContactResponse {
	var phone string
	setting, err := uc.settingRepo.Get("contact_whatsapp")
	if err == nil && setting != nil {
		phone = setting.Value
	}

	var email string
	emailSetting, err := uc.settingRepo.Get("contact_email")
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

func (uc *activationUsecaseImpl) logAttempt(licenseID *uuid.UUID, attemptedCode, hardwareID, ipAddress, status string) {
	log := domain.NewActivationLog(licenseID, attemptedCode, hardwareID, ipAddress, status)
	_ = uc.activationLogRepo.Create(log)
}
