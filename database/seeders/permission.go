package main

import (
	"ipincamp/srikandi-sehat/src/models"
	"log"

	"gorm.io/gorm"
)

func SeedPermissions(tx *gorm.DB) error {
	permissions := []models.Permission{
		// {Name: "view-profile"},
		// {Name: "edit-profile"},
	}

	for _, permission := range permissions {
		if err := tx.FirstOrCreate(&permission, models.Permission{Name: permission.Name}).Error; err != nil {
			return err
		}
	}

	log.Println("Permissions seeded")
	return nil
}
