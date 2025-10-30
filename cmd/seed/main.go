package main

import (
	"ipincamp/srikandi-sehat/config"
	"ipincamp/srikandi-sehat/database"
	"ipincamp/srikandi-sehat/database/seeders"
	"log"

	"gorm.io/gorm"
)

func main() {
	config.LoadConfig()
	config.SetTimeZone()
	database.ConnectDB()
	log.Println("[DB] [SEED] Starting seeding process...")

	tx := database.DB.Begin()
	if tx.Error != nil {
		log.Fatalf("[DB] [SEED] Failed to start transaction: %v", tx.Error)
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Fatalf("[DB] [SEED] Panic recovered, rolling back transaction: %v", r)
		}
	}()

	if err := seeds(tx); err != nil {
		tx.Rollback()
		log.Fatalf("[DB] [SEED] Seeding failed: %v", err)
	}

	if err := tx.Commit().Error; err != nil {
		log.Fatalf("[DB] [SEED] Failed to commit transaction: %v", err)
	}

	log.Println("[DB] [SEED] All seeding processes completed successfully!")
}

func seeds(tx *gorm.DB) error {
	if err := seeders.SeedRoles(tx); err != nil {
		return err
	}
	if err := seeders.SeedPermissions(tx); err != nil {
		return err
	}
	if err := seeders.SeedUsers(tx); err != nil {
		return err
	}
	if err := seeders.SeedRegions(tx); err != nil {
		return err
	}
	if err := seeders.SeedMenstrualData(tx); err != nil {
		return err
	}
	if config.Get("APP_ENV") == "production" {
		log.Println("[DB] [SEED] Skipping simulation data seeding in production environment.")
	} else {
		if err := seeders.SeedSimulationData(tx); err != nil {
			return err
		}
	}
	// And more...
	return nil
}
