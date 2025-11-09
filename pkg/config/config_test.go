package config_test

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ipincamp/srikandi-sehat/pkg/config"
)

// createTestEnvFile membuat file .env sementara untuk pengujian.
func createTestEnvFile(t *testing.T) {
	content := `
APP_PORT=9090
APP_HOST=127.0.0.1
APP_ENV=staging
APP_TIMEZONE=UTC

DB_HOST=localhost-test
DB_PORT=5433
DB_NAME=test_db
DB_USER=test_user
DB_PASS=test_pass
DB_SSL_MODE=disable
DB_TIMEZONE=UTC

TOKEN_ISSUER="test-issuer"
TOKEN_ACCESS_DURATION="1h"
TOKEN_REFRESH_DURATION="24h"
`
	// Kita simpan sebagai .env.staging agar tidak bentrok
	err := os.WriteFile(".env.staging", []byte(content), 0644)
	require.NoError(t, err)
}

func TestLoadConfig_Success(t *testing.T) {
	// 1. Setup: Buat file .env.staging
	createTestEnvFile(t)

	// 2. Setup: Ganti nama file yang dicari viper
	// (Viper tidak punya cara mudah untuk "ganti nama file", jadi kita set env var)
	// Trik: Kita akan memuat dari file yang namanya ".env"
	// Kita buat file .env baru untuk tes ini
	content, _ := os.ReadFile(".env.staging")
	_ = os.WriteFile(".env", content, 0644)

	// 3. Eksekusi
	cfg, err := config.LoadConfig(".") // Cari di direktori saat ini

	// 4. Validasi
	require.NoError(t, err)
	require.NotNil(t, cfg)

	assert.Equal(t, "9090", cfg.AppPort)
	assert.Equal(t, "localhost-test", cfg.DBHost)
	assert.Equal(t, "test_db", cfg.DBName)
	assert.Equal(t, "test_user", cfg.DBUsername)
	assert.Equal(t, "test-issuer", cfg.TokenIssuer)
	assert.Equal(t, 1*time.Hour, cfg.TokenAccessDuration)
	assert.Equal(t, 24*time.Hour, cfg.TokenRefreshDuration)
	assert.Equal(t, "staging", cfg.AppEnv)

	// 5. Cleanup
	os.Remove(".env.staging")
	os.Remove(".env")
}

func TestLoadConfig_Failure_Validation(t *testing.T) {
	// 1. Setup: Buat file .env yang tidak lengkap
	content := `
APP_PORT=9090
# APP_HOST hilang
APP_ENV=staging
APP_TIMEZONE=UTC

DB_HOST=localhost-test
DB_PORT=5433
DB_NAME=test_db
DB_USER=test_user
DB_PASS=test_pass
DB_SSL_MODE=disable
DB_TIMEZONE=UTC

TOKEN_ISSUER="test-issuer"
TOKEN_ACCESS_DURATION="1h"
TOKEN_REFRESH_DURATION="24h"
`
	err := os.WriteFile(".env", []byte(content), 0644)
	require.NoError(t, err)

	// 2. Eksekusi
	cfg, err := config.LoadConfig(".")

	// 3. Validasi
	require.Error(t, err)
	assert.Nil(t, cfg)
	// Cek pesan error spesifik dari validator
	assert.Contains(t, err.Error(), "field 'AppHost'")
	assert.Contains(t, err.Error(), "tidak memenuhi syarat 'required'")

	// 4. Cleanup
	os.Remove(".env")
}

func TestLoadConfig_Failure_InvalidEnvValue(t *testing.T) {
	// 1. Setup: Buat file .env dengan nilai enum yang salah
	content := `
APP_PORT=9090
APP_HOST=127.0.0.1
APP_ENV=test # 'test' tidak ada di 'oneof=development production staging'
APP_TIMEZONE=UTC

DB_HOST=localhost-test
DB_PORT=5433
DB_NAME=test_db
DB_USER=test_user
DB_PASS=test_pass
DB_SSL_MODE=disable
DB_TIMEZONE=UTC

TOKEN_ISSUER="test-issuer"
TOKEN_ACCESS_DURATION="1h"
TOKEN_REFRESH_DURATION="24h"
`
	err := os.WriteFile(".env", []byte(content), 0644)
	require.NoError(t, err)

	// 2. Eksekusi
	cfg, err := config.LoadConfig(".")

	// 3. Validasi
	require.Error(t, err)
	assert.Nil(t, cfg)
	assert.Contains(t, err.Error(), "field 'AppEnv'")
	assert.Contains(t, err.Error(), "tidak memenuhi syarat 'oneof'")

	// 4. Cleanup
	os.Remove(".env")
}
