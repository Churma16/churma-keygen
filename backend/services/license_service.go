package services

import (
	"crypto/rand"
	"errors"
	"fmt"

	"churma-keygen/backend/dtos"
	"churma-keygen/backend/models"
	"churma-keygen/backend/repositories"

	"github.com/google/uuid"
)

type LicenseService interface {
	GetAll() ([]dtos.LicenseResponse, error)
	Generate(req dtos.GenerateLicenseRequest) (*dtos.LicenseResponse, error)
	UpdateStatus(id string, req dtos.UpdateLicenseStatusRequest) (*dtos.LicenseResponse, error)
	Delete(id string) error
	GetActivationLogs() ([]dtos.ActivationLogResponse, error)
}

type licenseServiceImpl struct {
	licenseRepo       repositories.LicenseRepository
	clientRepo        repositories.ClientRepository
	activationLogRepo repositories.ActivationLogRepository
}

func NewLicenseService(
	licenseRepo repositories.LicenseRepository,
	clientRepo repositories.ClientRepository,
	activationLogRepo repositories.ActivationLogRepository,
) LicenseService {
	return &licenseServiceImpl{
		licenseRepo:       licenseRepo,
		clientRepo:        clientRepo,
		activationLogRepo: activationLogRepo,
	}
}

func (s *licenseServiceImpl) GetAll() ([]dtos.LicenseResponse, error) {
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

func (s *licenseServiceImpl) Generate(req dtos.GenerateLicenseRequest) (*dtos.LicenseResponse, error) {
	clientUUID, err := uuid.Parse(req.ClientID)
	if err != nil {
		return nil, errors.New("invalid client ID format")
	}

	// Verify client exists
	client, err := s.clientRepo.FindByID(clientUUID)
	if err != nil {
		return nil, errors.New("client not found")
	}

	trialLimit := req.TrialLimit
	if trialLimit < 0 {
		trialLimit = 100
	}

	licenseCode := generateRandomCode()

	license := models.License{
		ID:          uuid.New(),
		ClientID:    clientUUID,
		LicenseCode: licenseCode,
		TrialLimit:  trialLimit,
		Status:      "UNASSIGNED",
		ExpiresAt:   req.ExpiresAt,
	}

	err = s.licenseRepo.Create(&license)
	if err != nil {
		return nil, errors.New("failed to generate license key")
	}

	// Attach preloaded client
	license.Client = client

	res := mapLicenseToResponse(license)
	return &res, nil
}

func (s *licenseServiceImpl) UpdateStatus(id string, req dtos.UpdateLicenseStatusRequest) (*dtos.LicenseResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid license ID format")
	}

	// Validate status
	validStatuses := map[string]bool{
		"UNASSIGNED": true,
		"ACTIVE":     true,
		"SUSPENDED":  true,
		"REVOKED":    true,
	}
	if !validStatuses[req.Status] {
		return nil, errors.New("invalid status. Must be UNASSIGNED, ACTIVE, SUSPENDED, or REVOKED")
	}

	license, err := s.licenseRepo.FindByID(uid)
	if err != nil {
		return nil, errors.New("license not found")
	}

	license.Status = req.Status
	if req.Status == "UNASSIGNED" {
		license.HardwareID = ""
		license.ActivatedAt = nil
	}

	err = s.licenseRepo.Update(license)
	if err != nil {
		return nil, err
	}

	res := mapLicenseToResponse(*license)
	return &res, nil
}

func (s *licenseServiceImpl) Delete(id string) error {
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

func (s *licenseServiceImpl) GetActivationLogs() ([]dtos.ActivationLogResponse, error) {
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

func generateRandomCode() string {
	const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, 12)
	bytes := make([]byte, 12)
	_, _ = rand.Read(bytes)
	for i, b := range bytes {
		result[i] = chars[b%byte(len(chars))]
	}
	return fmt.Sprintf("SFA-%s-%s-%s", string(result[0:4]), string(result[4:8]), string(result[8:12]))
}

func mapLicenseToResponse(l models.License) dtos.LicenseResponse {
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

func mapActivationLogToResponse(al models.ActivationLog) dtos.ActivationLogResponse {
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
