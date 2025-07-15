package seeders

import (
	"car-service/constant"
	"car-service/database"
	"car-service/models"
	"log"
	"os"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// SeedersToAdmin seeds an admin if not already present
func SeedersToAdmin() {
	adminEmail := os.Getenv("ADMIN_EMAIL")
	adminPassword := os.Getenv("ADMIN_PASSWORD")

	var admin models.Register

	// Check if admin already exists
	err := database.DB.Where("email = ?", adminEmail).First(&admin).Error
	if err == nil {
		log.Println("✅ Admin already exists. Skipping seeding.")
		return
	}

	if err != nil && err != gorm.ErrRecordNotFound {
		log.Fatalf(" Error checking admin existence: %v", err)
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(adminPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf(" Error hashing password: %v", err)
	}

	// Create new admin
	seedAdmin := models.Register{
		Email:    adminEmail,
		Password: string(hashedPassword),
		Role:     constant.Admin,
		FirstName: "Admin",
	}

	if err := database.DB.Create(&seedAdmin).Error; err != nil {
		log.Fatalf(" Failed to create admin: %v", err)
	}

	log.Println(" Admin created successfully.")
}
