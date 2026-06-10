package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"churma-keygen/backend/config"
	"churma-keygen/backend/controllers"
	"churma-keygen/backend/crypto"
	"churma-keygen/backend/dtos"
	"churma-keygen/backend/middleware"
	"churma-keygen/backend/models"
	"churma-keygen/backend/repositories"
	"churma-keygen/backend/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	// Load environment variables in non-production environments
	if os.Getenv("ENV") != "prod" {
		err := godotenv.Load()
		if err != nil {
			fmt.Println("Warning: .env file not found, using system environment variables.")
		}
	}

	// 1. Connect to Database
	config.ConnectDatabase()
	db := config.DB

	// 2. Run GORM Auto Migrations
	fmt.Println("Running database migrations...")
	err := db.AutoMigrate(
		&models.User{},
		&models.Client{},
		&models.License{},
		&models.ActivationLog{},
	)
	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}
	fmt.Println("Database migrations completed successfully.")

	// 3. Instantiate Repositories
	userRepo := repositories.NewUserRepository(db)
	clientRepo := repositories.NewClientRepository(db)
	licenseRepo := repositories.NewLicenseRepository(db)
	activationLogRepo := repositories.NewActivationLogRepository(db)

	// 4. Seed Default Admin User if empty
	var adminCount int64
	db.Model(&models.User{}).Where("username = ?", "admin").Count(&adminCount)
	if adminCount == 0 {
		fmt.Println("Seeding default administrator...")
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
		if err != nil {
			log.Fatalf("Failed to hash default admin password: %v", err)
		}
		
		admin := models.User{
			ID:           uuid.New(),
			Username:     "admin",
			PasswordHash: string(hashedPassword),
			Role:         "SUPERADMIN",
		}
		if err := userRepo.Create(&admin); err != nil {
			log.Fatalf("Failed to seed default admin: %v", err)
		}
		fmt.Println("Default administrator created (username: admin, password: admin123).")
	}

	// 5. Initialize RSA Keypair
	err = crypto.InitRSAKeys()
	if err != nil {
		log.Fatalf("Failed to initialize RSA keys: %v", err)
	}

	// 6. Instantiate Services
	authService := services.NewAuthService(userRepo)
	clientService := services.NewClientService(clientRepo, licenseRepo)
	licenseService := services.NewLicenseService(licenseRepo, clientRepo, activationLogRepo)
	activationService := services.NewActivationService(licenseRepo, activationLogRepo)

	// 7. Instantiate Controllers
	authCtrl := controllers.NewAuthController(authService)
	clientCtrl := controllers.NewClientController(clientService)
	licenseCtrl := controllers.NewLicenseController(licenseService)
	activationCtrl := controllers.NewActivationController(activationService)

	// 8. Setup Gin Web Server
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
	}

	// Protected Admin Dashboard Routes
	adminGroup := r.Group("/api/v1/admin")
	adminGroup.Use(middleware.AdminAuth())
	{
		// Profile info
		adminGroup.GET("/me", authCtrl.GetMe)

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
	}

	// Serve Svelte frontend SPA built assets statically if the dist/ directory is present
	if _, err := os.Stat("./dist"); err == nil {
		fmt.Println("Serving Svelte production build from ./dist...")
		r.Static("/assets", "./dist/assets")
		r.NoRoute(func(c *gin.Context) {
			if len(c.Request.URL.Path) >= 4 && c.Request.URL.Path[:4] == "/api" {
				c.JSON(http.StatusNotFound, dtos.NewErrorResponse(http.StatusNotFound, "API route not found"))
				return
			}
			c.File("./dist/index.html")
		})
	}

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	address := fmt.Sprintf(":%s", port)
	
	fmt.Printf("Server listening on http://localhost:%s\n", port)
	if err := r.Run(address); err != nil {
		log.Fatalf("Failed to run web server: %v", err)
	}
}
