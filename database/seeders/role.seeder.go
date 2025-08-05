package seeders

import (
	"ipincamp/srikandi-sehat/src/constants"
	"ipincamp/srikandi-sehat/src/models"
	"log"

	"gorm.io/gorm"
)

func SeedRoles(tx *gorm.DB) error {
	log.Println("[DB] [SEED] [ROLE] Seeding roles...")
	roles := []models.Role{
		{Name: string(constants.AdminRole)},
		{Name: string(constants.UserRole)},
	}

	for _, role := range roles {
		if err := tx.FirstOrCreate(&role, models.Role{Name: role.Name}).Error; err != nil {
			return err
		}
	}

	log.Println("[DB] [SEED] [ROLE] Roles seeded successfully.")
	return nil
}
