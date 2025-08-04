package main

import (
	"ipincamp/srikandi-sehat/config"
	"ipincamp/srikandi-sehat/database"
	"ipincamp/srikandi-sehat/database/seeders"
	"log"
	"time"

	"gorm.io/gorm"
)

func main() {
	jakartaTime, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		log.Fatalf("[DB] [SEED] Failed to load Jakarta timezone: %v", err)
	}
	time.Local = jakartaTime
	config.LoadConfig()
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
	if err := seeders.SeedUsers(tx); err != nil {
		return err
	}
	// And more...
	return nil
}
