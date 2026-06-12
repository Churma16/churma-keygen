package services

import (
	"churma-keygen/backend/config"
	"churma-keygen/backend/dtos"
	"churma-keygen/backend/models"
	"churma-keygen/backend/repositories"
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func TestLicenseService_Generate_Unlimited(t *testing.T) {
	// Setup in-memory sqlite
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open SQLite in-memory: %v", err)
	}
	config.DB = db

	err = db.AutoMigrate(&models.Client{}, &models.License{}, &models.ActivationLog{})
	if err != nil {
		t.Fatalf("Failed to migrate models: %v", err)
	}

	clientRepo := repositories.NewClientRepository(db)
	licenseRepo := repositories.NewLicenseRepository(db)
	logRepo := repositories.NewActivationLogRepository(db)

	service := NewLicenseService(licenseRepo, clientRepo, logRepo)

	// Create a test client first
	client := models.Client{
		ID:        uuid.New(),
		Name:      "Test Client",
		OwnerName: "Test Owner",
		Phone:     "0812345",
	}
	if err := db.Create(&client).Error; err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// 1. Generate license key with TrialLimit = -1 (Unlimited)
	req := dtos.GenerateLicenseRequest{
		ClientID:   client.ID.String(),
		TrialLimit: -1,
	}

	resp, err := service.Generate(req)
	if err != nil {
		t.Fatalf("Failed to generate license: %v", err)
	}

	if resp.TrialLimit != -1 {
		t.Errorf("Expected TrialLimit to be -1, got %d", resp.TrialLimit)
	}

	// Retrieve from database to verify GORM didn't replace it
	var dbLicense models.License
	if err := db.First(&dbLicense, "id = ?", resp.ID).Error; err != nil {
		t.Fatalf("Failed to retrieve license from DB: %v", err)
	}

	if dbLicense.TrialLimit != -1 {
		t.Errorf("Expected TrialLimit in DB to be -1, got %d", dbLicense.TrialLimit)
	}

	// 2. Generate license key with TrialLimit = -2 (Invalid negative value, should be replaced with default 100)
	req2 := dtos.GenerateLicenseRequest{
		ClientID:   client.ID.String(),
		TrialLimit: -2,
	}

	resp2, err := service.Generate(req2)
	if err != nil {
		t.Fatalf("Failed to generate license 2: %v", err)
	}

	if resp2.TrialLimit != 100 {
		t.Errorf("Expected TrialLimit to be 100, got %d", resp2.TrialLimit)
	}
}
