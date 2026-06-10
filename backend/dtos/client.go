package dtos

import "time"

type CreateClientRequest struct {
	Name      string `json:"name" binding:"required"`
	OwnerName string `json:"owner_name"`
	Phone     string `json:"phone"`
}

type UpdateClientRequest struct {
	Name      string `json:"name" binding:"required"`
	OwnerName string `json:"owner_name"`
	Phone     string `json:"phone"`
}

type ClientResponse struct {
	ID        string            `json:"id"`
	Name      string            `json:"name"`
	OwnerName string            `json:"owner_name"`
	Phone     string            `json:"phone"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
	Licenses  []LicenseResponse `json:"licenses,omitempty"`
}

type ClientStatsResponse struct {
	TotalClients       int64 `json:"total_clients"`
	ActiveLicenses     int64 `json:"active_licenses"`
	SuspendedLicenses  int64 `json:"suspended_licenses"`
	UnassignedLicenses int64 `json:"unassigned_licenses"`
	RevokedLicenses    int64 `json:"revoked_licenses"`
}
