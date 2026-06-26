package main

import (
	"fmt"
	"log"
	"os"

	"churma-keygen/backend/config"
	"churma-keygen/backend/controllers"
	"churma-keygen/backend/crypto"
	"churma-keygen/backend/domain"
	"churma-keygen/backend/repositories"
	"churma-keygen/backend/usecase"

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
		&domain.User{},
		&domain.Client{},
		&domain.License{},
		&domain.ActivationLog{},
		&domain.Setting{},
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
	settingRepo := repositories.NewSettingRepository(db)

	// 4. Seed Default Admin User if empty
	var adminCount int64
	db.Model(&domain.User{}).Where("username = ?", "admin").Count(&adminCount)
	if adminCount == 0 {
		fmt.Println("Seeding default administrator...")
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
		if err != nil {
			log.Fatalf("Failed to hash default admin password: %v", err)
		}

		admin := domain.User{
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

	// Seed Default Setting if empty
	var settingCount int64
	db.Model(&domain.Setting{}).Where("key = ?", "contact_whatsapp").Count(&settingCount)
	if settingCount == 0 {
		fmt.Println("Seeding default WhatsApp contact...")
		defaultSetting := domain.Setting{
			Key:   "contact_whatsapp",
			Value: "",
		}
		if err := db.Create(&defaultSetting).Error; err != nil {
			log.Fatalf("Failed to seed default WhatsApp contact: %v", err)
		}
		fmt.Println("Default WhatsApp contact seeded.")
	}

	var emailSettingCount int64
	db.Model(&domain.Setting{}).Where("key = ?", "contact_email").Count(&emailSettingCount)
	if emailSettingCount == 0 {
		fmt.Println("Seeding default Email contact...")
		defaultEmailSetting := domain.Setting{
			Key:   "contact_email",
			Value: "",
		}
		if err := db.Create(&defaultEmailSetting).Error; err != nil {
			log.Fatalf("Failed to seed default Email contact: %v", err)
		}
		fmt.Println("Default Email contact seeded.")
	}

	// 5. Initialize RSA Keypair
	err = crypto.InitRSAKeys()
	if err != nil {
		log.Fatalf("Failed to initialize RSA keys: %v", err)
	}

	// 6. Instantiate Usecases
	authUsecase := usecase.NewAuthUsecase(userRepo)
	clientUsecase := usecase.NewClientUsecase(clientRepo, licenseRepo)
	licenseUsecase := usecase.NewLicenseUsecase(licenseRepo, clientRepo, activationLogRepo)
	activationUsecase := usecase.NewActivationUsecase(licenseRepo, activationLogRepo, settingRepo)
	settingUsecase := usecase.NewSettingUsecase(settingRepo)

	// 7. Instantiate Controllers
	authCtrl := controllers.NewAuthController(authUsecase)
	clientCtrl := controllers.NewClientController(clientUsecase)
	licenseCtrl := controllers.NewLicenseController(licenseUsecase)
	activationCtrl := controllers.NewActivationController(activationUsecase)
	settingCtrl := controllers.NewSettingController(settingUsecase)

	// 8. Setup Router
	r := SetupRouter(authCtrl, clientCtrl, licenseCtrl, activationCtrl, settingCtrl)

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
