package main

import (
	"net/http"
	"os"

	"churma-keygen/backend/controllers"
	"churma-keygen/backend/dtos"
	"churma-keygen/backend/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter(
	authCtrl *controllers.AuthController,
	clientCtrl *controllers.ClientController,
	licenseCtrl *controllers.LicenseController,
	activationCtrl *controllers.ActivationController,
	settingCtrl *controllers.SettingController,
) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	if os.Getenv("ENV") != "prod" {
		gin.SetMode(gin.DebugMode)
	}

	r := gin.Default()

	// Apply CORS
	r.Use(middleware.CORSMiddleware())

	// Health Check using struct response
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, dtos.NewSuccessResponse(http.StatusOK, "Churma Keygen API is active and running!", nil))
	})

	// Public Routes
	api := r.Group("/api/v1")
	{
		// Administrator Login
		api.POST("/auth/login", authCtrl.Login)

		// Public Client Activation Gateway (with Rate Limiting)
		api.POST("/client/activate", middleware.RateLimiter(), activationCtrl.ActivateLicense)

		// Public Endpoint to retrieve the RSA Public Key
		api.GET("/public-key", activationCtrl.GetPublicKey)

		// Public Endpoint to retrieve the support/contact details (WhatsApp)
		api.GET("/contact", activationCtrl.GetContact)
	}

	// Protected Admin Dashboard Routes
	adminGroup := r.Group("/api/v1/admin")
	adminGroup.Use(middleware.AdminAuth())
	{
		// Profile info
		adminGroup.GET("/me", authCtrl.GetMe)
		adminGroup.PUT("/profile", authCtrl.UpdateProfile)

		// Clients Management
		adminGroup.GET("/clients", clientCtrl.GetClients)
		adminGroup.POST("/clients", clientCtrl.CreateClient)
		adminGroup.PUT("/clients/:id", clientCtrl.UpdateClient)
		adminGroup.DELETE("/clients/:id", clientCtrl.DeleteClient)

		// Stats
		adminGroup.GET("/stats", clientCtrl.GetClientStats)

		// Licenses Management
		adminGroup.GET("/licenses", licenseCtrl.GetLicenses)
		adminGroup.POST("/licenses", licenseCtrl.GenerateLicense)
		adminGroup.PUT("/licenses/:id/status", licenseCtrl.UpdateLicenseStatus)
		adminGroup.DELETE("/licenses/:id", licenseCtrl.DeleteLicense)

		// Audit Activation Logs
		adminGroup.GET("/logs", licenseCtrl.GetActivationLogs)

		// System Settings Management
		adminGroup.GET("/settings/:key", settingCtrl.GetSetting)
		adminGroup.PUT("/settings/:key", settingCtrl.UpdateSetting)
	}

	// Serve Svelte frontend SPA built assets statically if the dist/ directory is present
	if _, err := os.Stat("./dist"); err == nil {
		r.Static("/assets", "./dist/assets")
		r.NoRoute(func(c *gin.Context) {
			if len(c.Request.URL.Path) >= 4 && c.Request.URL.Path[:4] == "/api" {
				c.JSON(http.StatusNotFound, dtos.NewErrorResponse(http.StatusNotFound, "API route not found"))
				return
			}
			c.File("./dist/index.html")
		})
	}

	return r
}
