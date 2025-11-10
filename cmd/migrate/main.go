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

// dropAllTables finds and drops all tables in the public schema using a safe method.
func dropAllTables(db *gorm.DB, logger zerolog.Logger) error {
	logger.Info().Msg("Finding all tables to drop...")

	// Get all table names in the 'public' schema
	// We don't need to filter out the migrations table, as we'll drop all.
	query := `SELECT table_name FROM information_schema.tables WHERE table_schema = 'public' AND table_type = 'BASE TABLE'`

	rows, err := db.Raw(query).Rows()
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
		logger.Info().Msg("No tables found to drop.")
		return nil
	}

	logger.Info().Msgf("Disabling foreign key constraints and dropping %d tables...", len(tables))

	// Start a transaction to manage the process
	return db.Transaction(func(tx *gorm.DB) error {
		// Disable Foreign Key constraints (PostgreSQL specific)
		// This allows us to drop tables in any order without constraint errors.
		if err := tx.Exec("SET session_replication_role = 'replica'").Error; err != nil {
			return fmt.Errorf("failed to disable foreign keys: %w", err)
		}

		// Drop all tables using GORM's migrator (safe from injection)
		for _, table := range tables {
			logger.Debug().Str("table", table).Msg("Dropping table...")
			// Use Migrator().DropTable() as recommended
			// This handles quoting and is safe from SQL injection.
			if err := tx.Migrator().DropTable(table); err != nil {
				// We log the error but attempt to continue,
				// as some failures might be expected if order is complex.
				logger.Warn().Err(err).Str("table", table).Msg("Failed to drop table, continuing...")
			}
		}

		// Re-enable Foreign Key constraints
		if err := tx.Exec("SET session_replication_role = 'origin'").Error; err != nil {
			return fmt.Errorf("failed to re-enable foreign keys: %w", err)
		}

		logger.Info().Int("count", len(tables)).Msg("All tables dropped successfully.")
		return nil
	})
}
