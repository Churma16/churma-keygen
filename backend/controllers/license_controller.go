package controllers

import (
	"crypto/rand"
	"fmt"
	"net/http"
	"time"

	"churma-keygen/backend/config"
	"churma-keygen/backend/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type GenerateLicenseRequest struct {
	ClientID   string     `json:"client_id" binding:"required"`
	TrialLimit int        `json:"trial_limit"` // defaults to 100
	ExpiresAt  *time.Time `json:"expires_at"`  // optional
}

type UpdateLicenseStatusRequest struct {
	Status string `json:"status" binding:"required"` // ACTIVE, SUSPENDED, REVOKED, UNASSIGNED
}

func GetLicenses(c *gin.Context) {
	var licenses []models.License
	err := config.DB.Preload("Client").Order("created_at DESC").Find(&licenses).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch licenses"})
		return
	}

	c.JSON(http.StatusOK, licenses)
}

func GenerateLicense(c *gin.Context) {
	var req GenerateLicenseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	clientUUID, err := uuid.Parse(req.ClientID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid client ID format"})
		return
	}

	// Verify client exists
	var client models.Client
	if err := config.DB.First(&client, "id = ?", clientUUID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Client not found"})
		return
	}

	trialLimit := req.TrialLimit
	if trialLimit < 0 {
		trialLimit = 100 // default
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

	if err := config.DB.Create(&license).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate license key"})
		return
	}

	// Preload client to return the complete object
	config.DB.Preload("Client").First(&license, "id = ?", license.ID)

	c.JSON(http.StatusCreated, license)
}

func UpdateLicenseStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid license ID format"})
		return
	}

	var req UpdateLicenseStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate status
	validStatuses := map[string]bool{
		"UNASSIGNED": true,
		"ACTIVE":     true,
		"SUSPENDED":  true,
		"REVOKED":    true,
	}
	if !validStatuses[req.Status] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status. Must be UNASSIGNED, ACTIVE, SUSPENDED, or REVOKED"})
		return
	}

	var license models.License
	if err := config.DB.First(&license, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "License not found"})
		return
	}

	license.Status = req.Status
	// If state changes, adjust timestamp accordingly or clear hardware_id if resetting
	if req.Status == "UNASSIGNED" {
		license.HardwareID = ""
		license.ActivatedAt = nil
	}

	if err := config.DB.Save(&license).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update license status"})
		return
	}

	c.JSON(http.StatusOK, license)
}

func DeleteLicense(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid license ID format"})
		return
	}

	var license models.License
	if err := config.DB.First(&license, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "License not found"})
		return
	}

	if err := config.DB.Delete(&license).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete license"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "License deleted successfully"})
}

func GetActivationLogs(c *gin.Context) {
	var logs []models.ActivationLog
	// Preload license and its client
	err := config.DB.Preload("License").Preload("License.Client").Order("created_at DESC").Limit(100).Find(&logs).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch activation logs"})
		return
	}

	c.JSON(http.StatusOK, logs)
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
