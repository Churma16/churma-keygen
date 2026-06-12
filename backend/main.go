package main

import (
	"fmt"
	"log"
	"os"

	"churma-keygen/backend/config"
	"churma-keygen/backend/controllers"
	"churma-keygen/backend/crypto"
	"churma-keygen/backend/models"
	"churma-keygen/backend/repositories"
	"churma-keygen/backend/services"

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

	// 8. Setup Router
	r := SetupRouter(authCtrl, clientCtrl, licenseCtrl, activationCtrl)

	// 9. Start server
	host := os.Getenv("HOST")
	if host == "" {
		host = "localhost"
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	address := fmt.Sprintf("%s:%s", host, port)

	fmt.Printf("Server listening on http://%s\n", address)
	if err := r.Run(address); err != nil {
		log.Fatalf("Failed to run web server: %v", err)
	}
}
