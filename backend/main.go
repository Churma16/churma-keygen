package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"churma-keygen/backend/config"
	"churma-keygen/backend/controllers"
	"churma-keygen/backend/crypto"
	"churma-keygen/backend/middleware"
	"churma-keygen/backend/models"

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

	// 3. Seed Default Admin User if empty
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
		if err := db.Create(&admin).Error; err != nil {
			log.Fatalf("Failed to seed default admin: %v", err)
		}
		fmt.Println("Default administrator created (username: admin, password: admin123).")
	}

	// 4. Initialize RSA Keypair
	err = crypto.InitRSAKeys()
	if err != nil {
		log.Fatalf("Failed to initialize RSA keys: %v", err)
	}

	// 5. Setup Gin Web Server
	gin.SetMode(gin.ReleaseMode)
	if os.Getenv("ENV") != "prod" {
		gin.SetMode(gin.DebugMode)
	}

	r := gin.Default()

	// Apply CORS
	r.Use(middleware.CORSMiddleware())

	// Health Check
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Churma Keygen API is active and running!"})
	})

	// Public Routes
	api := r.Group("/api/v1")
	{
		// Administrator Login
		api.POST("/auth/login", controllers.Login)

		// Public Client Activation Gateway (with Rate Limiting)
		api.POST("/client/activate", middleware.RateLimiter(), controllers.ActivateLicense)

		// Public Endpoint to retrieve the RSA Public Key (for clients to hardcode/verify JWT)
		api.GET("/public-key", controllers.GetPublicKey)
	}

	// Protected Admin Dashboard Routes
	admin := r.Group("/api/v1/admin")
	admin.Use(middleware.AdminAuth())
	{
		// Profile info
		admin.GET("/me", controllers.GetMe)

		// Clients Management
		admin.GET("/clients", controllers.GetClients)
		admin.POST("/clients", controllers.CreateClient)
		admin.PUT("/clients/:id", controllers.UpdateClient)
		admin.DELETE("/clients/:id", controllers.DeleteClient)

		// Stats
		admin.GET("/stats", controllers.GetClientStats)

		// Licenses Management
		admin.GET("/licenses", controllers.GetLicenses)
		admin.POST("/licenses", controllers.GenerateLicense)
		admin.PUT("/licenses/:id/status", controllers.UpdateLicenseStatus)
		admin.DELETE("/licenses/:id", controllers.DeleteLicense)

		// Audit Activation Logs
		admin.GET("/logs", controllers.GetActivationLogs)
	}

	// Serve Svelte frontend SPA built assets statically if the dist/ directory is present
	if _, err := os.Stat("./dist"); err == nil {
		fmt.Println("Serving Svelte production build from ./dist...")
		r.Static("/assets", "./dist/assets")
		r.NoRoute(func(c *gin.Context) {
			if len(c.Request.URL.Path) >= 4 && c.Request.URL.Path[:4] == "/api" {
				c.JSON(http.StatusNotFound, gin.H{"error": "API route not found"})
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
