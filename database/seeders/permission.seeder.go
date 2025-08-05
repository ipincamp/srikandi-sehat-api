package seeders

import (
	"ipincamp/srikandi-sehat/src/models"
	"log"

	"gorm.io/gorm"
)

func SeedPermissions(tx *gorm.DB) error {
	log.Println("[DB] [SEED] [PERMISSION] Seeding permissions...")
	permissions := []models.Permission{
		// {Name: "permission-name"},
	}

	for _, permission := range permissions {
		if err := tx.FirstOrCreate(&permission, models.Permission{Name: permission.Name}).Error; err != nil {
			return err
		}
	}

	log.Println("[DB] [SEED] [PERMISSION] Permissions seeded successfully.")
	return nil
}
