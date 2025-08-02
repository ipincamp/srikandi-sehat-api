package main

import (
	"ipincamp/srikandi-sehat/config"
	"ipincamp/srikandi-sehat/database"
	"ipincamp/srikandi-sehat/database/seeders/seed"
	"log"

	"gorm.io/gorm"
)

func main() {
	config.LoadConfig()
	database.ConnectDB()

	tx := database.DB.Begin()
	if tx.Error != nil {
		log.Fatalf("Failed to start transaction: %v", tx.Error)
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Fatalf("Panic recovered, rolling back transaction: %v", r)
		}
	}()

	if err := SeedAll(tx); err != nil {
		tx.Rollback()
		log.Fatalf("Seeding failed: %v", err)
	}

	if err := tx.Commit().Error; err != nil {
		log.Fatalf("Failed to commit transaction: %v", err)
	}

	log.Println("All seeding processes completed successfully!")
}

func SeedAll(db *gorm.DB) error {
	if err := seed.SeedRoles(db); err != nil {
		return err
	}
	if err := seed.SeedPermissions(db); err != nil {
		return err
	}
	if err := seed.SeedUsers(db); err != nil {
		return err
	}
	if err := seed.SeedRegions(db); err != nil {
		return err
	}

	return nil
}
