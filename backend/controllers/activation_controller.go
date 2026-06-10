package controllers

import (
	"errors"
	"net/http"
	"time"

	"churma-keygen/backend/config"
	"churma-keygen/backend/crypto"
	"churma-keygen/backend/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ActivateRequest struct {
	LicenseCode string `json:"license_code" binding:"required"`
	HardwareID  string `json:"hardware_id" binding:"required"`
}

func ActivateLicense(c *gin.Context) {
	var req ActivateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "License code and hardware ID are required"})
		return
	}

	ipAddress := c.ClientIP()

	// 1. Find the license
	var license models.License
	err := config.DB.Preload("Client").Where("license_code = ?", req.LicenseCode).First(&license).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Log as INVALID_KEY
			logAttempt(nil, req.LicenseCode, req.HardwareID, ipAddress, "INVALID_KEY")
			c.JSON(http.StatusNotFound, gin.H{"error": "Invalid license code"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// 2. Check if license is suspended or revoked
	if license.Status == "SUSPENDED" || license.Status == "REVOKED" {
		logAttempt(&license.ID, req.LicenseCode, req.HardwareID, ipAddress, "SUSPENDED_KEY")
		c.JSON(http.StatusForbidden, gin.H{"error": "This license key is suspended or revoked"})
		return
	}

	// 3. Check if license is expired
	if license.ExpiresAt != nil && license.ExpiresAt.Before(time.Now()) {
		logAttempt(&license.ID, req.LicenseCode, req.HardwareID, ipAddress, "SUSPENDED_KEY")
		c.JSON(http.StatusForbidden, gin.H{"error": "This license key has expired"})
		return
	}

	// 4. Handle Hardware ID binding
	if license.HardwareID == "" {
		// First-time activation: bind the HWID
		license.HardwareID = req.HardwareID
		license.Status = "ACTIVE"
		now := time.Now()
		license.ActivatedAt = &now

		if err := config.DB.Save(&license).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to activate license"})
			return
		}

		logAttempt(&license.ID, req.LicenseCode, req.HardwareID, ipAddress, "SUCCESS")
	} else {
		// Subsequent activation: verify HWID
		if license.HardwareID != req.HardwareID {
			logAttempt(&license.ID, req.LicenseCode, req.HardwareID, ipAddress, "HWID_MISMATCH")
			c.JSON(http.StatusConflict, gin.H{"error": "Hardware ID mismatch. This license is bound to another machine."})
			return
		}

		logAttempt(&license.ID, req.LicenseCode, req.HardwareID, ipAddress, "SUCCESS")
	}

	// 5. Generate signed JWT token
	token, err := crypto.SignLicenseToken(
		license.ClientID.String(),
		license.Client.Name,
		license.HardwareID,
		license.TrialLimit,
		license.ExpiresAt,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to sign license token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token":       token,
		"status":      license.Status,
		"expires_at":  license.ExpiresAt,
		"trial_limit": license.TrialLimit,
		"client_name": license.Client.Name,
	})
}

// GetPublicKey endpoint to serve the RSA public key to client apps
func GetPublicKey(c *gin.Context) {
	pubPEM, err := crypto.GetPublicKeyPEM()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "RSA Public Key is not configured"})
		return
	}

	// Return PEM format as plain text or JSON
	c.String(http.StatusOK, pubPEM)
}

func logAttempt(licenseID *uuid.UUID, attemptedCode, hardwareID, ipAddress, status string) {
	log := models.ActivationLog{
		ID:                uuid.New(),
		LicenseID:         licenseID,
		AttemptedCode:     attemptedCode,
		HardwareIDAttempt: hardwareID,
		IPAddress:         ipAddress,
		Status:            status,
		CreatedAt:         time.Now(),
	}
	// We run db insert in a new transaction, ignoring error so we don't block request
	_ = config.DB.Create(&log).Error
}
