package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"churma-keygen/backend/config"
	"churma-keygen/backend/crypto"
	"churma-keygen/backend/dtos"
	"churma-keygen/backend/models"
	"churma-keygen/backend/repositories"
	"churma-keygen/backend/services"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	testClient  models.Client
	testLicense models.License
)

func setupTestEnvironment(t *testing.T) *gin.Engine {
	// 1. Setup in-memory SQLite database
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open test SQLite: %v", err)
	}
	config.DB = db

	// Migrate models
	err = db.AutoMigrate(
		&models.User{},
		&models.Client{},
		&models.License{},
		&models.ActivationLog{},
		&models.Setting{},
	)
	if err != nil {
		t.Fatalf("Failed to migrate test models: %v", err)
	}

	// 2. Generate RSA Keypair for signing
	err = crypto.GenerateRSAKeyPair("test_private.pem", "test_public.pem")
	if err != nil {
		t.Fatalf("Failed to generate test RSA keypair: %v", err)
	}
	os.Setenv("RSA_PRIVATE_KEY_PATH", "test_private.pem")
	os.Setenv("RSA_PUBLIC_KEY_PATH", "test_public.pem")

	err = crypto.InitRSAKeys()
	if err != nil {
		t.Fatalf("Failed to load test keys: %v", err)
	}

	// 3. Seed test data
	testClient = models.Client{
		ID:        uuid.New(),
		Name:      "Toko Test Sejahtera",
		OwnerName: "Budi",
		Phone:     "08121212",
	}
	if err := db.Create(&testClient).Error; err != nil {
		t.Fatalf("Failed to seed client: %v", err)
	}

	testLicense = models.License{
		ID:          uuid.New(),
		ClientID:    testClient.ID,
		LicenseCode: "TEST-XXXX-YYYY",
		TrialLimit:  50,
		Status:      "UNASSIGNED",
	}
	if err := db.Create(&testLicense).Error; err != nil {
		t.Fatalf("Failed to seed license: %v", err)
	}

	testSetting := models.Setting{
		Key:   "contact_whatsapp",
		Value: "0812-3456-7890",
	}
	if err := db.Create(&testSetting).Error; err != nil {
		t.Fatalf("Failed to seed test setting: %v", err)
	}

	// 4. Instantiate Layers for Testing
	licenseRepo := repositories.NewLicenseRepository(db)
	activationLogRepo := repositories.NewActivationLogRepository(db)
	settingRepo := repositories.NewSettingRepository(db)
	activationService := services.NewActivationService(licenseRepo, activationLogRepo, settingRepo)
	activationCtrl := NewActivationController(activationService)

	// 5. Setup router
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/api/v1/client/activate", activationCtrl.ActivateLicense)
	r.GET("/api/v1/contact", activationCtrl.GetContact)
	return r
}

func teardownTestEnvironment() {
	os.Remove("test_private.pem")
	os.Remove("test_public.pem")
	os.Unsetenv("RSA_PRIVATE_KEY_PATH")
	os.Unsetenv("RSA_PUBLIC_KEY_PATH")
}

func TestActivateLicense_Success_FirstTime(t *testing.T) {
	r := setupTestEnvironment(t)
	defer teardownTestEnvironment()

	reqBody := dtos.ActivateRequest{
		LicenseCode: "TEST-XXXX-YYYY",
		HardwareID:  "HWID-COMPUTER-1",
	}
	bodyBytes, _ := json.Marshal(reqBody)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/client/activate", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d. Response: %s", w.Code, w.Body.String())
	}

	var envelope map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &envelope)

	meta, ok := envelope["meta"].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected meta object in response")
	}
	if meta["status"] != "success" {
		t.Errorf("Expected meta status to be success, got %v", meta["status"])
	}

	data, ok := envelope["data"].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected data object in response")
	}

	if data["status"] != "ACTIVE" {
		t.Errorf("Expected data status to be ACTIVE, got %v", data["status"])
	}
	if data["client_name"] != "Toko Test Sejahtera" {
		t.Errorf("Expected client name Toko Test Sejahtera, got %v", data["client_name"])
	}
	if data["token"] == nil || data["token"] == "" {
		t.Errorf("Expected JWT token to be issued")
	}

	// Verify database changes
	var lic models.License
	config.DB.First(&lic, "id = ?", testLicense.ID)
	if lic.HardwareID != "HWID-COMPUTER-1" {
		t.Errorf("Expected HWID to be bound as HWID-COMPUTER-1, got %s", lic.HardwareID)
	}
	if lic.Status != "ACTIVE" {
		t.Errorf("Expected database license status to be ACTIVE, got %s", lic.Status)
	}

	// Verify audit log
	var logAttempt models.ActivationLog
	err := config.DB.Where("license_id = ?", testLicense.ID).First(&logAttempt).Error
	if err != nil {
		t.Errorf("Expected activation log to be written: %v", err)
	}
	if logAttempt.Status != "SUCCESS" {
		t.Errorf("Expected log status SUCCESS, got %s", logAttempt.Status)
	}
}

func TestActivateLicense_Success_Reactivation(t *testing.T) {
	r := setupTestEnvironment(t)
	defer teardownTestEnvironment()

	// Pretend it's already active with HWID-COMPUTER-1
	config.DB.Model(&testLicense).Updates(map[string]interface{}{
		"status":      "ACTIVE",
		"hardware_id": "HWID-COMPUTER-1",
	})

	reqBody := dtos.ActivateRequest{
		LicenseCode: "TEST-XXXX-YYYY",
		HardwareID:  "HWID-COMPUTER-1",
	}
	bodyBytes, _ := json.Marshal(reqBody)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/client/activate", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200 on reactivation, got %d. Response: %s", w.Code, w.Body.String())
	}
}

func TestActivateLicense_Error_HWIDMismatch(t *testing.T) {
	r := setupTestEnvironment(t)
	defer teardownTestEnvironment()

	// Active on MACHINE-A
	config.DB.Model(&testLicense).Updates(map[string]interface{}{
		"status":      "ACTIVE",
		"hardware_id": "MACHINE-A",
	})

	reqBody := dtos.ActivateRequest{
		LicenseCode: "TEST-XXXX-YYYY",
		HardwareID:  "MACHINE-B", // trying to spoof on different computer
	}
	bodyBytes, _ := json.Marshal(reqBody)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/client/activate", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)

	if w.Code != http.StatusConflict {
		t.Errorf("Expected status 409 Conflict on HWID mismatch, got %d. Response: %s", w.Code, w.Body.String())
	}

	// Check mismatch log
	var logAttempt models.ActivationLog
	config.DB.Order("created_at desc").First(&logAttempt)
	if logAttempt.Status != "HWID_MISMATCH" {
		t.Errorf("Expected log status HWID_MISMATCH, got %s", logAttempt.Status)
	}
}

func TestActivateLicense_Error_InvalidCode(t *testing.T) {
	r := setupTestEnvironment(t)
	defer teardownTestEnvironment()

	reqBody := dtos.ActivateRequest{
		LicenseCode: "WRONG-CODE-1234",
		HardwareID:  "HWID-SOME-MACHINE",
	}
	bodyBytes, _ := json.Marshal(reqBody)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/client/activate", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404 on invalid code, got %d", w.Code)
	}

	// Verify log recorded INVALID_KEY
	var logAttempt models.ActivationLog
	config.DB.Order("created_at desc").First(&logAttempt)
	if logAttempt.Status != "INVALID_KEY" {
		t.Errorf("Expected log status INVALID_KEY, got %s", logAttempt.Status)
	}
}

func TestActivateLicense_Error_Suspended(t *testing.T) {
	r := setupTestEnvironment(t)
	defer teardownTestEnvironment()

	// Suspended status
	config.DB.Model(&testLicense).Update("status", "SUSPENDED")

	reqBody := dtos.ActivateRequest{
		LicenseCode: "TEST-XXXX-YYYY",
		HardwareID:  "HWID-COMPUTER-1",
	}
	bodyBytes, _ := json.Marshal(reqBody)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/client/activate", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("Expected status 403 Forbidden on suspended key, got %d", w.Code)
	}

	// Verify log recorded SUSPENDED_KEY
	var logAttempt models.ActivationLog
	config.DB.Order("created_at desc").First(&logAttempt)
	if logAttempt.Status != "SUSPENDED_KEY" {
		t.Errorf("Expected log status SUSPENDED_KEY, got %s", logAttempt.Status)
	}
}

func TestActivateLicense_Error_Expired(t *testing.T) {
	r := setupTestEnvironment(t)
	defer teardownTestEnvironment()

	// Set expired time in the past
	pastTime := time.Now().Add(-24 * time.Hour)
	config.DB.Model(&testLicense).Update("expires_at", &pastTime)

	reqBody := dtos.ActivateRequest{
		LicenseCode: "TEST-XXXX-YYYY",
		HardwareID:  "HWID-COMPUTER-1",
	}
	bodyBytes, _ := json.Marshal(reqBody)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/client/activate", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("Expected status 403 Forbidden on expired key, got %d", w.Code)
	}

	// Verify log recorded SUSPENDED_KEY
	var logAttempt models.ActivationLog
	config.DB.Order("created_at desc").First(&logAttempt)
	if logAttempt.Status != "SUSPENDED_KEY" {
		t.Errorf("Expected log status SUSPENDED_KEY, got %s", logAttempt.Status)
	}
}

func TestGetContact_Success(t *testing.T) {
	r := setupTestEnvironment(t)
	defer teardownTestEnvironment()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/contact", nil)

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d. Response: %s", w.Code, w.Body.String())
	}

	var envelope map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &envelope)

	meta, ok := envelope["meta"].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected meta object in response")
	}
	if meta["status"] != "success" {
		t.Errorf("Expected meta status to be success, got %v", meta["status"])
	}

	data, ok := envelope["data"].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected data object in response")
	}

	// In setupTestEnvironment, we seeded contact_whatsapp = "0812-3456-7890"
	// After sanitization, phone is still "0812-3456-7890", but URL is https://wa.me/6281234567890
	if data["phone"] != "0812-3456-7890" {
		t.Errorf("Expected phone to be 0812-3456-7890, got %v", data["phone"])
	}
	if data["whatsapp_url"] != "https://wa.me/6281234567890" {
		t.Errorf("Expected whatsapp_url to be https://wa.me/6281234567890, got %v", data["whatsapp_url"])
	}
}
