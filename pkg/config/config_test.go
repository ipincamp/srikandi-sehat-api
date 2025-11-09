package config_test

import (
	"strings"
	"testing"

	"github.com/ipincamp/srikandi-sehat/pkg/config"
)

// TestLoad verifies the Load function in various scenarios.
func TestLoad(t *testing.T) {

	// --- Sub-test 1: Success Case ---
	// Tests if config is loaded correctly when all env vars are provided.
	t.Run("SuccessCase_AllVariablesSet", func(t *testing.T) {
		// Use t.Setenv to safely set environment variables for this sub-test.
		// These variables are automatically restored after the test.
		t.Setenv("PORT", "9090")
		t.Setenv("ENV", "staging")
		t.Setenv("DB_HOST", "db.staging.com")
		t.Setenv("DB_PORT", "5433")
		t.Setenv("DB_USER", "staging_user")
		t.Setenv("DB_PASS", "staging_pass")
		t.Setenv("DB_NAME", "staging_db")
		t.Setenv("DB_SSL_MODE", "require")

		cfg, err := config.Load()

		// 1. Check for unexpected errors
		if err != nil {
			t.Fatalf("Load() returned an unexpected error: %v", err)
		}
		if cfg == nil {
			t.Fatal("Load() returned nil config on success")
		}

		// 2. Check Server config
		if cfg.Server.Port != "9090" {
			t.Errorf("got Server.Port %q, want %q", cfg.Server.Port, "9090")
		}
		if cfg.Server.Env != "staging" {
			t.Errorf("got Server.Env %q, want %q", cfg.Server.Env, "staging")
		}

		// 3. Check Database config
		if cfg.Database.Host != "db.staging.com" {
			t.Errorf("got Database.Host %q, want %q", cfg.Database.Host, "db.staging.com")
		}
		if cfg.Database.Port != 5433 {
			t.Errorf("got Database.Port %d, want %d", cfg.Database.Port, 5433)
		}
		if cfg.Database.User != "staging_user" {
			t.Errorf("got Database.User %q, want %q", cfg.Database.User, "staging_user")
		}
		// Note: We don't test password output for security, just that it's loaded.
		if cfg.Database.Password != "staging_pass" {
			t.Errorf("got Database.Password %q, want %q", cfg.Database.Password, "staging_pass")
		}
		if cfg.Database.DBName != "staging_db" {
			t.Errorf("got Database.DBName %q, want %q", cfg.Database.DBName, "staging_db")
		}
		if cfg.Database.SSLMode != "require" {
			t.Errorf("got Database.SSLMode %q, want %q", cfg.Database.SSLMode, "require")
		}
	})

	// --- Sub-test 2: Default Values ---
	// Tests if default values are applied when optional env vars are missing.
	t.Run("DefaultValues_WhenOptionalVarsAreMissing", func(t *testing.T) {
		// We *must* set the required variables to pass validation.
		t.Setenv("DB_USER", "required_user")
		t.Setenv("DB_PASS", "required_pass")
		t.Setenv("DB_NAME", "required_db")

		// We are *not* setting PORT, ENV, DB_HOST, DB_PORT, or DB_SSL_MODE
		// to verify their default fallback values.

		cfg, err := config.Load()
		if err != nil {
			t.Fatalf("Load() returned an unexpected error: %v", err)
		}

		// Check default values
		if cfg.Server.Port != "8080" {
			t.Errorf("got default Server.Port %q, want %q", cfg.Server.Port, "8080")
		}
		if cfg.Server.Env != "development" {
			t.Errorf("got default Server.Env %q, want %q", cfg.Server.Env, "development")
		}
		if cfg.Database.Host != "localhost" {
			t.Errorf("got default Database.Host %q, want %q", cfg.Database.Host, "localhost")
		}
		if cfg.Database.Port != 5432 {
			t.Errorf("got default Database.Port %d, want %d", cfg.Database.Port, 5432)
		}
		if cfg.Database.SSLMode != "disable" {
			t.Errorf("got default Database.SSLMode %q, want %q", cfg.Database.SSLMode, "disable")
		}
	})

	// --- Sub-test 3: Validation Error ---
	// Tests if Load() fails when required variables are missing.
	t.Run("ValidationError_MissingRequiredVars", func(t *testing.T) {
		// In this test, we set *no* variables.
		// `t.Setenv` ensures a clean slate, so DB_USER, DB_PASS,
		// and DB_NAME will be empty strings.

		cfg, err := config.Load()

		if err == nil {
			t.Fatal("Load() did not return an error, but was expected to")
		}
		if cfg != nil {
			t.Error("Load() returned a non-nil config on error")
		}

		// Check for the specific validation error message
		expectedErr := "DB_USER, DB_PASS, and DB_NAME must be set"
		if err.Error() != expectedErr {
			t.Errorf("got error %q, want %q", err.Error(), expectedErr)
		}
	})

	// --- Sub-test 4: Parse Error ---
	// Tests if Load() fails when DB_PORT is not a valid integer.
	t.Run("ParseError_InvalidDbPort", func(t *testing.T) {
		// Set required fields to pass validation
		t.Setenv("DB_USER", "user")
		t.Setenv("DB_PASS", "pass")
		t.Setenv("DB_NAME", "db")
		// Set the invalid port
		t.Setenv("DB_PORT", "not-a-number")

		_, err := config.Load()

		if err == nil {
			t.Fatal("Load() did not return an error for invalid DB_PORT, but was expected to")
		}

		// Check for the specific parse error prefix
		if !strings.HasPrefix(err.Error(), "invalid DB_PORT:") {
			t.Errorf("got error %q, want prefix 'invalid DB_PORT:'", err.Error())
		}
	})
}

// TestDSN verifies the DSN method constructs the connection string correctly.
// This uses a table-driven test for multiple scenarios.
func TestDSN(t *testing.T) {
	// Define test cases
	testCases := []struct {
		name     string          // Name of the test case
		db       config.Database // Input
		expected string          // Expected output
	}{
		{
			name: "Full Production Config",
			db: config.Database{
				User:     "admin",
				Password: "password123",
				Host:     "my.host.com",
				Port:     5432,
				DBName:   "prod_db",
				SSLMode:  "verify-full",
			},
			expected: "postgres://admin:password123@my.host.com:5432/prod_db?sslmode=verify-full",
		},
		{
			name: "Localhost Development Config",
			db: config.Database{
				User:     "local_user",
				Password: "local_password",
				Host:     "localhost",
				Port:     5432,
				DBName:   "test_db",
				SSLMode:  "disable",
			},
			expected: "postgres://local_user:local_password@localhost:5432/test_db?sslmode=disable",
		},
		{
			name: "Config with Special Chars in Password",
			db: config.Database{
				User:     "user_with_chars",
				Password: "p@ss!w*rd#?",
				Host:     "127.0.0.1",
				Port:     5433,
				DBName:   "char_db",
				SSLMode:  "allow",
			},
			expected: "postgres://user_with_chars:p@ss!w*rd#?@127.0.0.1:5433/char_db?sslmode=allow",
		},
	}

	// Iterate over test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.db.DSN()
			if got != tc.expected {
				t.Errorf("DSN() mismatch:\ngot:  %s\nwant: %s", got, tc.expected)
			}
		})
	}
}
