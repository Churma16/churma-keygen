package dtos

import "time"

type ActivateRequest struct {
	LicenseCode string `json:"license_code" binding:"required"`
	HardwareID  string `json:"hardware_id" binding:"required"`
}

type ActivateResponse struct {
	Token      string     `json:"token"`
	Status     string     `json:"status"`
	ExpiresAt  *time.Time `json:"expires_at"`
	TrialLimit int        `json:"trial_limit"`
	ClientName string     `json:"client_name"`
}
