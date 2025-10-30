package main

import (
	"bufio"
	"fmt"
	"ipincamp/srikandi-sehat/config"
	"ipincamp/srikandi-sehat/database"
	"ipincamp/srikandi-sehat/database/migrations"
	"log"
	"os"
	"strings"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func dropAllTables(db *gorm.DB) error {
	log.Println("Dropping tables dynamically...")

	models := []any{
		"maintenance_whitelists",
		"notifications",
		"settings",
		"recommendations",
		"symptom_log_details",
		"symptom_logs",
		"symptom_options",
		"symptoms",
		"menstrual_cycles",
		"profiles",
		"villages",
		"districts",
		"regencies",
		"provinces",
		"classifications",
		"invalid_tokens",
		"user_roles",
		"user_permissions",
		"role_permissions",
		"roles",
		"permissions",
		"users",
	}

	if err := db.Migrator().DropTable(models...); err != nil {
		log.Printf("Error dropping GORM model tables: %v", err)
		return err
	}
	log.Println("Dynamically dropped GORM model tables.")

	migrationsTableName := gormigrate.DefaultOptions.TableName
	if db.Migrator().HasTable(migrationsTableName) {
		if err := db.Migrator().DropTable(migrationsTableName); err != nil {
			log.Printf("Could not drop migrations table '%s': %v", migrationsTableName, err)
			// Consider if this should be a fatal error depending on your needs
		} else {
			log.Printf("Dropped migrations table: %s", migrationsTableName)
		}
	} else {
		log.Printf("Migrations table '%s' does not exist, skipping drop.", migrationsTableName)
	}

	return nil
}

func askForConfirmation(prompt string) bool {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("%s [y/N]: ", prompt)
		response, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("Error reading confirmation: %v", err)
			return false
		}
		response = strings.ToLower(strings.TrimSpace(response))
		switch response {
		case "y", "yes":
			return true
		case "", "n", "no":
			return false
		}
	}
}

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
		migrations.AddLatePeriodNotifiedToMenstrualCycles(),
		migrations.CreateMaintenanceTables(),
		migrations.AddEmailVerification(),
		migrations.AddLastOtpSentAt(),
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
	case "reset":
		confirm := askForConfirmation("Are you sure you want to reset the database? This will drop all tables and data.")
		if !confirm {
			log.Println("Database reset cancelled.")
			return
		}
		if err := dropAllTables(db); err != nil {
			log.Fatalf("[DB] [MIGRATE] Failed to drop all tables: %v", err)
		}
		if err := m.Migrate(); err != nil {
			log.Fatalf("[DB] [MIGRATE] Failed to run migration after reset: %v", err)
		}
		log.Println("[DB] [MIGRATE] Database reset and migration completed successfully.")
	case "drop-all":
		confirm := askForConfirmation("Are you sure you want to drop all tables and data?")
		if !confirm {
			log.Println("Drop all tables cancelled.")
			return
		}
		if err := dropAllTables(db); err != nil {
			log.Fatalf("[DB] [MIGRATE] Failed to drop all tables: %v", err)
		}
		log.Println("[DB] [MIGRATE] Dropped all tables successfully.")
	default:
		log.Fatalf("[DB] [MIGRATE] Unknown command: %s", command)
	}
}
