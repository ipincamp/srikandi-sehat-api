package seeders

import (
	"ipincamp/srikandi-sehat/database/factories"
	"ipincamp/srikandi-sehat/src/constants"
	"ipincamp/srikandi-sehat/src/models"
	"log"

	"gorm.io/gorm"
)

func SeedUsers(tx *gorm.DB) error {
	log.Println("[DB] [SEED] [USER] Seeding admin user...")
	adminUser, err := factories.CreateAdminUser(tx)
	if err != nil {
		return err
	}
	var adminRole models.Role
	if err := tx.First(&adminRole, "name = ?", constants.AdminRole).Error; err != nil {
		return err
	}
	if err := tx.Model(&adminUser).Association("Roles").Replace(&adminRole); err != nil {
		return err
	}
	log.Println("[DB] [SEED] [USER] Admin user seeded successfully.")

	/*
		log.Println("[DB] [SEED] [USER] Creating 100 random users...")
		randomUsers, err := factories.CreateUsers(tx, 100)
		if err != nil {
			return err
		}
		log.Printf("[DB] [SEED] [USER] %d random users created successfully.", len(randomUsers))

		var userRole models.Role
		if err := tx.First(&userRole, "name = ?", constants.UserRole).Error; err != nil {
			return err
		}
		log.Println("[DB] [SEED] [USER] Assigning 'User' role to random users...")
		for _, user := range randomUsers {
			if err := tx.Model(&user).Association("Roles").Append(&userRole); err != nil {
				log.Printf("[DB] [SEED] [USER] Failed to assign role to user %s: %v", user.Email, err)
			}
		}
		log.Println("[DB] [SEED] [USER] Role assignment for random users completed.")

		log.Printf("[DB] [SEED] [USER] Seeding completed successfully with %d users.", len(randomUsers)+1)
	*/
	log.Print("[DB] [SEED] [USER] Seeding completed successfully.")
	return nil
}
