package postgres

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	"github.com/ipincamp/srikandi-sehat/pkg/config"
)

// ConnectGORM creates a new GORM database instance.
// This is used specifically by the migration tool.
func ConnectGORM(cfg config.Database) (*gorm.DB, error) {
	// Build the DSN from the loaded config
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
		cfg.Host,
		cfg.User,
		cfg.Password,
		cfg.DBName,
		cfg.Port,
		cfg.SSLMode,
		"Asia/Jakarta", // Hardcoded from .env, you might want to move this to config struct
	)

	// Configure GORM logger
	newLogger := gormlogger.New(
		log.New(log.Writer(), "\r\n", log.LstdFlags), // io writer
		gormlogger.Config{
			SlowThreshold:             time.Second,     // Slow SQL threshold
			LogLevel:                  gormlogger.Info, // Log level
			IgnoreRecordNotFoundError: true,            // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      true,            // Don't include params in output
			Colorful:                  false,           // Disable color
		},
	)

	// Open the database connection
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database using GORM: %w", err)
	}

	return db, nil
}
