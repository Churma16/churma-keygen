package dtos

import "time"

type GenerateLicenseRequest struct {
	ClientID   string     `json:"client_id" binding:"required"`
	TrialLimit int        `json:"trial_limit"`
	ExpiresAt  *time.Time `json:"expires_at"`
}

type UpdateLicenseStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

type LicenseResponse struct {
	ID          string     `json:"id"`
	ClientID    string     `json:"client_id"`
	ClientName  string     `json:"client_name,omitempty"`
	LicenseCode string     `json:"license_code"`
	HardwareID  string     `json:"hardware_id"`
	TrialLimit  int        `json:"trial_limit"`
	Status      string     `json:"status"`
	ExpiresAt   *time.Time `json:"expires_at"`
	ActivatedAt *time.Time `json:"activated_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type ActivationLogResponse struct {
	ID                 string     `json:"id"`
	LicenseID          *string    `json:"license_id,omitempty"`
	ClientName         string     `json:"client_name,omitempty"`
	AttemptedCode      string     `json:"attempted_code"`
	HardwareIDAttempt  string     `json:"hardware_id_attempt"`
	IPAddress          string     `json:"ip_address"`
	Status             string     `json:"status"`
	CreatedAt          time.Time  `json:"created_at"`
}
