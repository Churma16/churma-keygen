package usecase

import (
	"errors"

	"churma-keygen/backend/domain"
	"churma-keygen/backend/dtos"

	"github.com/google/uuid"
)

type LicenseUsecase interface {
	GetAll() ([]dtos.LicenseResponse, error)
	Generate(req dtos.GenerateLicenseRequest) (*dtos.LicenseResponse, error)
	UpdateStatus(id string, req dtos.UpdateLicenseStatusRequest) (*dtos.LicenseResponse, error)
	Delete(id string) error
	GetActivationLogs() ([]dtos.ActivationLogResponse, error)
}

type licenseUsecaseImpl struct {
	licenseRepo       domain.LicenseRepository
	clientRepo        domain.ClientRepository
	activationLogRepo domain.ActivationLogRepository
}

func NewLicenseUsecase(
	licenseRepo domain.LicenseRepository,
	clientRepo domain.ClientRepository,
	activationLogRepo domain.ActivationLogRepository,
) LicenseUsecase {
	return &licenseUsecaseImpl{
		licenseRepo:       licenseRepo,
		clientRepo:        clientRepo,
		activationLogRepo: activationLogRepo,
	}
}

func (s *licenseUsecaseImpl) GetAll() ([]dtos.LicenseResponse, error) {
	licenses, err := s.licenseRepo.FindAll()
	if err != nil {
		return nil, err
	}

	var resp []dtos.LicenseResponse
	for _, l := range licenses {
		resp = append(resp, mapLicenseToResponse(l))
	}
	return resp, nil
}

func (s *licenseUsecaseImpl) Generate(req dtos.GenerateLicenseRequest) (*dtos.LicenseResponse, error) {
	clientUUID, err := uuid.Parse(req.ClientID)
	if err != nil {
		return nil, errors.New("invalid client ID format")
	}

	client, err := s.clientRepo.FindByID(clientUUID)
	if err != nil {
		return nil, errors.New("client not found")
	}

	trialLimit := req.TrialLimit
	if trialLimit < -1 {
		trialLimit = 100
	}

	license := &domain.License{
		ID:         uuid.New(),
		ClientID:   clientUUID,
		TrialLimit: trialLimit,
		Status:     "UNASSIGNED",
		ExpiresAt:  req.ExpiresAt,
	}

	// Delegate to domain
	license.GenerateRandomCode()

	err = s.licenseRepo.Create(license)
	if err != nil {
		return nil, errors.New("failed to generate license key")
	}

	license.Client = client

	res := mapLicenseToResponse(*license)
	return &res, nil
}

func (s *licenseUsecaseImpl) UpdateStatus(id string, req dtos.UpdateLicenseStatusRequest) (*dtos.LicenseResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid license ID format")
	}

	license, err := s.licenseRepo.FindByID(uid)
	if err != nil {
		return nil, errors.New("license not found")
	}

	// Delegate logic to domain
	if err := license.UpdateStatus(req.Status); err != nil {
		return nil, err
	}

	err = s.licenseRepo.Update(license)
	if err != nil {
		return nil, err
	}

	res := mapLicenseToResponse(*license)
	return &res, nil
}

func (s *licenseUsecaseImpl) Delete(id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return errors.New("invalid license ID format")
	}

	_, err = s.licenseRepo.FindByID(uid)
	if err != nil {
		return errors.New("license not found")
	}

	return s.licenseRepo.Delete(uid)
}

func (s *licenseUsecaseImpl) GetActivationLogs() ([]dtos.ActivationLogResponse, error) {
	logs, err := s.activationLogRepo.FindAll(100)
	if err != nil {
		return nil, err
	}

	var resp []dtos.ActivationLogResponse
	for _, l := range logs {
		resp = append(resp, mapActivationLogToResponse(l))
	}
	return resp, nil
}

func mapLicenseToResponse(l domain.License) dtos.LicenseResponse {
	clientName := ""
	if l.Client != nil {
		clientName = l.Client.Name
	}
	return dtos.LicenseResponse{
		ID:          l.ID.String(),
		ClientID:    l.ClientID.String(),
		ClientName:  clientName,
		LicenseCode: l.LicenseCode,
		HardwareID:  l.HardwareID,
		TrialLimit:  l.TrialLimit,
		Status:      l.Status,
		ExpiresAt:   l.ExpiresAt,
		ActivatedAt: l.ActivatedAt,
		CreatedAt:   l.CreatedAt,
		UpdatedAt:   l.UpdatedAt,
	}
}

func mapActivationLogToResponse(al domain.ActivationLog) dtos.ActivationLogResponse {
	clientName := ""
	if al.License != nil && al.License.Client != nil {
		clientName = al.License.Client.Name
	}
	var licenseIDStr *string
	if al.LicenseID != nil {
		s := al.LicenseID.String()
		licenseIDStr = &s
	}
	return dtos.ActivationLogResponse{
		ID:                 al.ID.String(),
		LicenseID:          licenseIDStr,
		ClientName:         clientName,
		AttemptedCode:      al.AttemptedCode,
		HardwareIDAttempt:  al.HardwareIDAttempt,
		IPAddress:          al.IPAddress,
		Status:             al.Status,
		CreatedAt:          al.CreatedAt,
	}
}
