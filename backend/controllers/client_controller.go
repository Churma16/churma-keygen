package controllers

import (
	"net/http"

	"churma-keygen/backend/config"
	"churma-keygen/backend/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

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

func GetClients(c *gin.Context) {
	var clients []models.Client
	// Preload licenses to show them on dashboard
	err := config.DB.Preload("Licenses").Order("name ASC").Find(&clients).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch clients"})
		return
	}

	c.JSON(http.StatusOK, clients)
}

func GetClientStats(c *gin.Context) {
	var clientCount int64
	var activeCount int64
	var suspendedCount int64
	var unassignedCount int64
	var revokedCount int64

	config.DB.Model(&models.Client{}).Count(&clientCount)
	config.DB.Model(&models.License{}).Where("status = ?", "ACTIVE").Count(&activeCount)
	config.DB.Model(&models.License{}).Where("status = ?", "SUSPENDED").Count(&suspendedCount)
	config.DB.Model(&models.License{}).Where("status = ?", "UNASSIGNED").Count(&unassignedCount)
	config.DB.Model(&models.License{}).Where("status = ?", "REVOKED").Count(&revokedCount)

	c.JSON(http.StatusOK, gin.H{
		"total_clients":       clientCount,
		"active_licenses":     activeCount,
		"suspended_licenses":  suspendedCount,
		"unassigned_licenses": unassignedCount,
		"revoked_licenses":    revokedCount,
	})
}

func CreateClient(c *gin.Context) {
	var req CreateClientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	client := models.Client{
		ID:        uuid.New(),
		Name:      req.Name,
		OwnerName: req.OwnerName,
		Phone:     req.Phone,
	}

	err := config.DB.Create(&client).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create client"})
		return
	}

	c.JSON(http.StatusCreated, client)
}

func UpdateClient(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid client ID format"})
		return
	}

	var req UpdateClientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var client models.Client
	err = config.DB.First(&client, "id = ?", id).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Client not found"})
		return
	}

	client.Name = req.Name
	client.OwnerName = req.OwnerName
	client.Phone = req.Phone

	err = config.DB.Save(&client).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update client"})
		return
	}

	c.JSON(http.StatusOK, client)
}

func DeleteClient(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid client ID format"})
		return
	}

	var client models.Client
	err = config.DB.First(&client, "id = ?", id).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Client not found"})
		return
	}

	// GORM soft delete
	err = config.DB.Delete(&client).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete client"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Client deleted successfully"})
}
