package postgres

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	// maxRetries is the maximum number of retries for connecting to the database.
	maxRetries = 5
	// retryDelay is the duration to wait between retries.
	retryDelay = 5 * time.Second
)

// Connect establishes a connection pool to the PostgreSQL database.
// It takes a context (for cancellation) and a DSN (Data Source Name) string.
// It performs retries on connection failure, which is crucial for robust startup
// in containerized environments (e.g., waiting for the DB container).
//
// Returns a connection pool (*pgxpool.Pool) or an error if connection fails.
func Connect(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	var pool *pgxpool.Pool
	var err error

	// Parse the DSN to create a config.
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database DSN: %w", err)
	}

	// --- Connection Retry Logic ---
	for i := 0; i < maxRetries; i++ {
		// Try to connect
		pool, err = pgxpool.NewWithConfig(ctx, config)
		if err == nil {
			// Try to ping the database to verify the connection is live.
			if pingErr := pool.Ping(ctx); pingErr == nil {
				// Connection successful
				return pool, nil
			} else {
				// Ping failed, close the pool and log the error for this attempt
				pool.Close()
				err = pingErr
				log.Printf("DB Ping failed (attempt %d/%d): %v", i+1, maxRetries, err)
			}
		} else {
			log.Printf("DB Connect failed (attempt %d/%d): %v", i+1, maxRetries, err)
		}

		// Wait before retrying, unless it's the last attempt
		if i < maxRetries-1 {
			log.Printf("Waiting %s before next DB connection attempt...", retryDelay)
			time.Sleep(retryDelay)
		}
	}

	// If the loop finishes without returning, all retries have failed.
	return nil, fmt.Errorf("failed to connect to database after %d attempts: %w", maxRetries, err)
}
