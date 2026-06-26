package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"churma-keygen/backend/config"
	"churma-keygen/backend/dtos"
	"churma-keygen/backend/middleware"
	"churma-keygen/backend/models"
	"churma-keygen/backend/repositories"
	"churma-keygen/backend/usecase"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

func setupSettingTestEnvironment(t *testing.T) (*gin.Engine, *gorm.DB) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open SQLite: %v", err)
	}
	config.DB = db

	err = db.AutoMigrate(&models.Setting{})
	if err != nil {
		t.Fatalf("Failed to migrate settings model: %v", err)
	}

	settingRepo := repositories.NewSettingRepository(db)
	settingUsecase := usecase.NewSettingUsecase(settingRepo)
	settingCtrl := NewSettingController(settingUsecase)

	// Seed setting
	db.Create(&models.Setting{Key: "contact_whatsapp", Value: "6281234567890"})

	gin.SetMode(gin.TestMode)
	r := gin.New()

	// Admin routes group with auth
	adminGroup := r.Group("/api/v1/admin")
	adminGroup.Use(middleware.AdminAuth())
	{
		adminGroup.GET("/settings/:key", settingCtrl.GetSetting)
		adminGroup.PUT("/settings/:key", settingCtrl.UpdateSetting)
	}

	return r, db
}

func generateTestToken(t *testing.T) string {
	claims := &middleware.AdminClaims{
		UserID:   "admin-id-123",
		Username: "admin",
		Role:     "SUPERADMIN",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "churma-keygen-admin",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(middleware.JWTSecret)
	if err != nil {
		t.Fatalf("Failed to sign test token: %v", err)
	}
	return tokenString
}

func TestGetSetting_Success(t *testing.T) {
	r, _ := setupSettingTestEnvironment(t)
	token := generateTestToken(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/admin/settings/contact_whatsapp", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d. Response: %s", w.Code, w.Body.String())
	}

	var envelope map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &envelope)
	data, ok := envelope["data"].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected data object in response")
	}

	if data["value"] != "6281234567890" {
		t.Errorf("Expected setting value to be 6281234567890, got %v", data["value"])
	}
}

func TestUpdateSetting_Success(t *testing.T) {
	r, db := setupSettingTestEnvironment(t)
	token := generateTestToken(t)

	reqBody := dtos.UpdateSettingRequest{
		Value: "628999999999",
	}
	bodyBytes, _ := json.Marshal(reqBody)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/v1/admin/settings/contact_whatsapp", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d. Response: %s", w.Code, w.Body.String())
	}

	// Verify database was updated
	var setting models.Setting
	db.Where("key = ?", "contact_whatsapp").First(&setting)
	if setting.Value != "628999999999" {
		t.Errorf("Expected database value to be updated to 628999999999, got %s", setting.Value)
	}
}

func TestGetSetting_Unauthorized(t *testing.T) {
	r, _ := setupSettingTestEnvironment(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/admin/settings/contact_whatsapp", nil)

	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
}
