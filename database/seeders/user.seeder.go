package seeders

import (
	"ipincamp/srikandi-sehat/config"
	"ipincamp/srikandi-sehat/database/factories"
	"ipincamp/srikandi-sehat/src/models"
	"ipincamp/srikandi-sehat/src/utils"
	"log"

	"gorm.io/gorm"
)

func SeedUsers(tx *gorm.DB) error {
	log.Println("[DB] [SEED] Seeding admin user...")
	adminPassword, _ := utils.HashPassword(config.Get("ADMIN_PASSWORD"))
	adminUser := models.User{
		Name:     config.Get("ADMIN_NAME"),
		Email:    config.Get("ADMIN_EMAIL"),
		Password: adminPassword,
	}
	if err := tx.Where(models.User{Email: adminUser.Email}).FirstOrCreate(&adminUser).Error; err != nil {
		return err
	}
	log.Printf("[DB] [SEED] Admin user created: %s", adminUser.Email)

	// TODO: Assign role to adminUser

	log.Println("[DB] [SEED] Creating 100 random users...")
	randomUsers, err := factories.CreateUsers(tx, 100)
	if err != nil {
		return err
	}
	log.Printf("[DB] [SEED] %d random users created successfully.", len(randomUsers))

	// TODO: Assign roles to random users

	return nil
}
