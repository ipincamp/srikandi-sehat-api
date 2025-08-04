package main

import (
	"ipincamp/srikandi-sehat/config"
	"ipincamp/srikandi-sehat/database"
	"ipincamp/srikandi-sehat/database/migrations"
	"log"
	"os"

	"github.com/go-gormigrate/gormigrate/v2"
)

func main() {
	config.LoadConfig()
	database.ConnectDB()
	db := database.DB

	m := gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		migrations.CreateUsersTable(),
		// And more...
	})

	if len(os.Args) < 2 {
		log.Fatal("[Database] Migration command is required (e.g., up, down)")
	}
	command := os.Args[1]

	switch command {
	case "up":
		if err := m.Migrate(); err != nil {
			log.Fatalf("[Database] Failed to run migration: %v", err)
		}
		log.Println("[Database] Migration completed successfully.")
	case "down":
		if err := m.RollbackLast(); err != nil {
			log.Fatalf("[Database] Failed to rollback migration: %v", err)
		}
		log.Println("[Database] Rollback last migration completed successfully.")
	default:
		log.Fatalf("[Database] Unknown command: %s", command)
	}
}
