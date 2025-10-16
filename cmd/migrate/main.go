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
		migrations.AddUuidAndSoftDeleteToUsers(),
		migrations.CreateRbacTables(),
		migrations.CreateInvalidTokensTable(),
		migrations.CreateRegionTables(),
		migrations.CreateProfilesTable(),
		migrations.CreateMenstrualTrackingTables(),
		migrations.AddIndexesToSymptomLogs(),
		migrations.AddCycleIdToSymptomLogs(),
		migrations.ChangeCycleDatesToDatetime(),
		migrations.RenameLogDateInSymptomLogs(),
		migrations.AddSoftDeleteToMenstrualCycles(),
		migrations.AddFcmTokenToUsersTable(),
		migrations.CreateNotificationsTable(),
		migrations.AddLongPeriodNotifiedToMenstrualCycles(),
		// And more...
	})

	if len(os.Args) < 2 {
		log.Fatal("[DB] [MIGRATE] Migration command is required (e.g., up, down)")
	}
	command := os.Args[1]

	switch command {
	case "up":
		if err := m.Migrate(); err != nil {
			log.Fatalf("[DB] [MIGRATE] Failed to run migration: %v", err)
		}
		log.Println("[DB] [MIGRATE] Migration completed successfully.")
	case "down":
		if err := m.RollbackLast(); err != nil {
			log.Fatalf("[DB] [MIGRATE] Failed to rollback migration: %v", err)
		}
		log.Println("[DB] [MIGRATE] Rollback last migration completed successfully.")
	default:
		log.Fatalf("[DB] [MIGRATE] Unknown command: %s", command)
	}
}
