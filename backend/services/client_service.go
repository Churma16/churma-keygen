package services

import (
	"errors"

	"churma-keygen/backend/dtos"
	"churma-keygen/backend/models"
	"churma-keygen/backend/repositories"

	"github.com/google/uuid"
)

type ClientService interface {
	GetAll() ([]dtos.ClientResponse, error)
	Create(req dtos.CreateClientRequest) (*dtos.ClientResponse, error)
	Update(id string, req dtos.UpdateClientRequest) (*dtos.ClientResponse, error)
	Delete(id string) error
	GetStats() (*dtos.ClientStatsResponse, error)
}

type clientServiceImpl struct {
	clientRepo  repositories.ClientRepository
	licenseRepo repositories.LicenseRepository
}

func NewClientService(clientRepo repositories.ClientRepository, licenseRepo repositories.LicenseRepository) ClientService {
	return &clientServiceImpl{
		clientRepo:  clientRepo,
		licenseRepo: licenseRepo,
	}
}

func (s *clientServiceImpl) GetAll() ([]dtos.ClientResponse, error) {
	clients, err := s.clientRepo.FindAll()
	if err != nil {
		return nil, err
	}

	var resp []dtos.ClientResponse
	for _, c := range clients {
		resp = append(resp, mapClientToResponse(c))
	}
	return resp, nil
}

func (s *clientServiceImpl) Create(req dtos.CreateClientRequest) (*dtos.ClientResponse, error) {
	client := models.Client{
		ID:        uuid.New(),
		Name:      req.Name,
		OwnerName: req.OwnerName,
		Phone:     req.Phone,
	}

	err := s.clientRepo.Create(&client)
	if err != nil {
		return nil, err
	}

	return &dtos.ClientResponse{
		ID:        client.ID.String(),
		Name:      client.Name,
		OwnerName: client.OwnerName,
		Phone:     client.Phone,
		CreatedAt: client.CreatedAt,
		UpdatedAt: client.UpdatedAt,
	}, nil
}

func (s *clientServiceImpl) Update(id string, req dtos.UpdateClientRequest) (*dtos.ClientResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid client ID format")
	}

	client, err := s.clientRepo.FindByID(uid)
	if err != nil {
		return nil, errors.New("client not found")
	}

	client.Name = req.Name
	client.OwnerName = req.OwnerName
	client.Phone = req.Phone

	err = s.clientRepo.Update(client)
	if err != nil {
		return nil, err
	}

	res := mapClientToResponse(*client)
	return &res, nil
}

func (s *clientServiceImpl) Delete(id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return errors.New("invalid client ID format")
	}

	_, err = s.clientRepo.FindByID(uid)
	if err != nil {
		return errors.New("client not found")
	}

	return s.clientRepo.Delete(uid)
}

func (s *clientServiceImpl) GetStats() (*dtos.ClientStatsResponse, error) {
	totalClients, err := s.clientRepo.Count()
	if err != nil {
		return nil, err
	}

	active, _ := s.licenseRepo.CountByStatus("ACTIVE")
	suspended, _ := s.licenseRepo.CountByStatus("SUSPENDED")
	unassigned, _ := s.licenseRepo.CountByStatus("UNASSIGNED")
	revoked, _ := s.licenseRepo.CountByStatus("REVOKED")

	return &dtos.ClientStatsResponse{
		TotalClients:       totalClients,
		ActiveLicenses:     active,
		SuspendedLicenses:  suspended,
		UnassignedLicenses: unassigned,
		RevokedLicenses:    revoked,
	}, nil
}

func mapClientToResponse(c models.Client) dtos.ClientResponse {
	var licenses []dtos.LicenseResponse
	for _, l := range c.Licenses {
		licenses = append(licenses, mapLicenseToResponse(l))
	}
	return dtos.ClientResponse{
		ID:        c.ID.String(),
		Name:      c.Name,
		OwnerName: c.OwnerName,
		Phone:     c.Phone,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
		Licenses:  licenses,
	}
}
