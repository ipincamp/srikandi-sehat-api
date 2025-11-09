package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	// Using godotenv for local development convenience.
	// In production, environment variables should be set directly.
	"github.com/joho/godotenv"
)

// Config holds all configuration for the application.
type Config struct {
	// Server holds web server specific configuration.
	Server Server
	// Database holds PostgreSQL database connection configuration.
	Database Database
	// Token holds PASETO token configuration.
	Token Token
}

// Server holds configuration related to the HTTP server.
type Server struct {
	// Port is the port number the server will listen on.
	Port string
	// Env is the application environment (e.g., "development", "staging", "production").
	Env string
}

// Database holds configuration for the PostgreSQL connection.
type Database struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// Token holds configuration for PASETO token generation.
type Token struct {
	SymmetricKey    string
	Issuer          string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

// Load reads configuration from environment variables.
// It loads from a .env file if it exists (for local development).
func Load() (*Config, error) {
	// Attempt to load .env file.
	// This is ignored if .env does not exist, which is fine for production.
	_ = godotenv.Load()

	dbPort, err := strconv.Atoi(getEnv("DB_PORT", "5432"))
	if err != nil {
		return nil, fmt.Errorf("invalid DB_PORT: %w", err)
	}

	accessTokenTTL, err := time.ParseDuration(getEnv("TOKEN_ACCESS_DURATION", "15m"))
	if err != nil {
		return nil, fmt.Errorf("invalid TOKEN_ACCESS_DURATION: %w", err)
	}

	refreshTokenTTL, err := time.ParseDuration(getEnv("TOKEN_REFRESH_DURATION", "720h"))
	if err != nil {
		return nil, fmt.Errorf("invalid TOKEN_REFRESH_DURATION: %w", err)
	}

	cfg := &Config{
		Server: Server{
			Port: getEnv("PORT", "8080"),
			Env:  getEnv("ENV", "development"),
		},
		Database: Database{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     dbPort,
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASS", ""),
			DBName:   getEnv("DB_NAME", "postgres"),
			SSLMode:  getEnv("DB_SSL_MODE", "disable"),
		},
		Token: Token{
			SymmetricKey:    getEnv("TOKEN_SYMMETRIC_KEY", ""),
			Issuer:          getEnv("TOKEN_ISSUER", "srikandi-sehat"),
			AccessTokenTTL:  accessTokenTTL,
			RefreshTokenTTL: refreshTokenTTL,
		},
	}

	// Simple validation
	if cfg.Database.User == "" || cfg.Database.Password == "" || cfg.Database.DBName == "" {
		return nil, fmt.Errorf("DB_USER, DB_PASS, and DB_NAME must be set")
	}
	if cfg.Token.SymmetricKey == "" {
		return nil, fmt.Errorf("TOKEN_SYMMETRIC_KEY must be set")
	}
	if len(cfg.Token.SymmetricKey) != 32 {
		return nil, fmt.Errorf("TOKEN_SYMMETRIC_KEY must be exactly 32 bytes")
	}

	return cfg, nil
}

// DSN returns the Data Source Name string for connecting to the database.
// This encapsulates the logic for building the connection string.
func (d *Database) DSN() string {
	// Example: "postgres://user:password@localhost:5432/dbname?sslmode=disable"
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		d.User,
		d.Password,
		d.Host,
		d.Port,
		d.DBName,
		d.SSLMode,
	)
}

// getEnv retrieves an environment variable by key or returns a fallback value.
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
