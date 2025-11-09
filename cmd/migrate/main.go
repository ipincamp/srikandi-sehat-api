package main

import (
	"fmt"
	"log"
	"os"

	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/rs/zerolog"
	"gorm.io/gorm"

	// Internal dependencies
	"github.com/ipincamp/srikandi-sehat/internal/adapters/driven/postgres"
	"github.com/ipincamp/srikandi-sehat/internal/adapters/driven/postgres/db/migrations"
	"github.com/ipincamp/srikandi-sehat/pkg/config"
	applogger "github.com/ipincamp/srikandi-sehat/pkg/logger"
)

func main() {
	// 1. Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// 2. Initialize logger (for standard output)
	// We use the app's logger for consistency
	logger := applogger.NewLogger(cfg.Server.Env)

	// 3. Connect to Database using GORM
	// This uses the new GORM connector, not the pgx one
	db, err := postgres.ConnectGORM(cfg.Database)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to connect to database for migration")
	}

	// 4. Initialize gormigrate
	m := gormigrate.New(db, gormigrate.DefaultOptions, migrations.AllMigrations)

	// 5. Get command-line argument
	if len(os.Args) < 2 {
		logger.Fatal().Msg("Missing command. Usage: go run ./cmd/migrate/main.go [up|down]")
	}
	command := os.Args[1]

	// 6. Execute command
	switch command {
	case "up":
		logger.Info().Msg("Applying migrations...")
		if err := m.Migrate(); err != nil {
			logger.Fatal().Err(err).Msg("Migration failed")
		}
		logger.Info().Msg("Migrations applied successfully")

	case "down":
		logger.Info().Msg("Rolling back last migration...")
		// RollbackLast reverts the last applied migration
		if err := m.RollbackLast(); err != nil {
			logger.Fatal().Err(err).Msg("Rollback failed")
		}
		logger.Info().Msg("Rollback successful")

	case "fresh":
		logger.Info().Msg("Resetting database...")
		// 1. Drop all tables
		if err := dropAllTables(db, logger); err != nil {
			logger.Fatal().Err(err).Msg("Failed to drop tables during reset")
		}
		// 2. Re-apply all migrations
		logger.Info().Msg("Applying all migrations...")
		if err := m.Migrate(); err != nil {
			logger.Fatal().Err(err).Msg("Migration failed during reset")
		}
		logger.Info().Msg("Database reset successfully")

	case "prune":
		logger.Info().Msg("Pruning database...")
		// 1. Just drop all tables
		if err := dropAllTables(db, logger); err != nil {
			logger.Fatal().Err(err).Msg("Failed to drop tables during prune")
		}
		logger.Info().Msg("Database pruned successfully (all tables dropped)")

	default:
		logger.Fatal().Msgf("Unknown command: '%s'. Use 'up', 'down', 'fresh', or 'prune'.", command)
	}
}

// dropAllTables finds and drops all tables in the public schema,
// except for the gormigrate history table.
func dropAllTables(db *gorm.DB, logger zerolog.Logger) error {
	logger.Info().Msg("Finding all tables to drop...")

	// Get all table names in the 'public' schema
	// We MUST exclude the migration history table, or gormigrate will fail
	query := `SELECT table_name FROM information_schema.tables WHERE table_schema = 'public' AND table_type = 'BASE TABLE' AND table_name != ?`

	rows, err := db.Raw(query, gormigrate.DefaultOptions.TableName).Rows()
	if err != nil {
		return fmt.Errorf("failed to query for tables: %w", err)
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return fmt.Errorf("failed to scan table name: %w", err)
		}
		tables = append(tables, name)
	}

	if len(tables) == 0 {
		logger.Info().Msg("No user tables found to drop.")
		return nil
	}

	logger.Info().Msgf("Found %d tables to drop. Dropping...", len(tables))

	// Drop tables one by one. Using CASCADE handles foreign key constraints.
	for _, table := range tables {
		logger.Debug().Str("table", table).Msg("Dropping table...")
		// We must quote the table name to handle any special characters or casing
		if err := db.Exec(fmt.Sprintf(`DROP TABLE IF EXISTS "public"."%s" CASCADE`, table)).Error; err != nil {
			return fmt.Errorf("failed to drop table %s: %w", table, err)
		}
	}

	if e := db.Exec(`DROP TABLE IF EXISTS "public"."migrations" CASCADE`).Error; e != nil {
		return fmt.Errorf("failed to drop table migrations: %w", e)
	}

	logger.Info().Int("count", len(tables)).Msg("All user tables dropped successfully.")
	return nil
}
