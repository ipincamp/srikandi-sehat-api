package config

import (
	"errors"
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

// Config adalah struct yang menampung semua konfigurasi aplikasi.
// Kita menggunakan tag 'validate' untuk memastikan nilai-nilai penting ada.
type Config struct {
	// Server
	AppEnv  string `mapstructure:"APP_ENV"      validate:"required,oneof=development production staging"`
	AppTZ   string `mapstructure:"APP_TIMEZONE" validate:"required"`
	AppPort string `mapstructure:"APP_PORT"     validate:"required"`
	AppHost string `mapstructure:"APP_HOST"     validate:"required"`

	// Database
	DBHost     string `mapstructure:"DB_HOST"     validate:"required"`
	DBPort     string `mapstructure:"DB_PORT"     validate:"required"`
	DBName     string `mapstructure:"DB_NAME"     validate:"required"`
	DBUsername string `mapstructure:"DB_USER"     validate:"required"`
	DBPassword string `mapstructure:"DB_PASS"     validate:"required"`
	DBSSLMode  string `mapstructure:"DB_SSL_MODE" validate:"required"`
	DBTimezone string `mapstructure:"DB_TIMEZONE" validate:"required"`

	// Paseto
	TokenIssuer          string        `mapstructure:"TOKEN_ISSUER"           validate:"required"`
	TokenAccessDuration  time.Duration `mapstructure:"TOKEN_ACCESS_DURATION"  validate:"required"`
	TokenRefreshDuration time.Duration `mapstructure:"TOKEN_REFRESH_DURATION" validate:"required"`
}

// validate adalah instance dari validator.
var validate = validator.New()

// LoadConfig memuat konfigurasi dari file .env di path yang diberikan.
// Ini juga akan membaca dari environment variables (berguna untuk Docker/K8s).
func LoadConfig(path string) (*Config, error) {
	// 1. Set default, file, dan env binding
	viper.AddConfigPath(path)   // Path untuk mencari file config
	viper.SetConfigName(".env") // Nama file config (tanpa ekstensi)
	viper.SetConfigType("env")  // Tipe file config

	viper.AutomaticEnv() // Baca env variables yang cocok

	// 2. Baca file config
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// File config tidak ditemukan; tidak masalah jika env vars di-set
		} else {
			return nil, fmt.Errorf("gagal membaca file config: %w", err)
		}
	}

	// 3. Unmarshal config ke struct
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("gagal unmarshal config: %w", err)
	}

	// 4. Validasi struct
	if err := validate.Struct(&config); err != nil {
		// Ubah error validasi menjadi lebih mudah dibaca
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			for _, e := range validationErrors {
				return nil, fmt.Errorf("validasi config gagal: field '%s' (value: '%v') tidak memenuhi syarat '%s'", e.Field(), e.Value(), e.Tag())
			}
		}
		return nil, fmt.Errorf("validasi config gagal: %w", err)
	}

	return &config, nil
}
