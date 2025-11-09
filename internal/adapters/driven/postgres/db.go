package postgres

import (
	"fmt"
	"log"

	"github.com/ipincamp/srikandi-sehat/pkg/config"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

// ConnectToDB adalah constructor untuk adapter database.
// Ini mengimplementasikan "Driven Adapter" [cite: 78] dan
// mengambil konfigurasi sebagai dependensi.
func ConnectToDB(cfg *config.Config) (*sqlx.DB, error) {
	// 1. Buat Data Source Name (DSN) string dari config
	// dsn := "host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai"
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBUsername,
		cfg.DBPassword,
		cfg.DBName,
		cfg.DBSSLMode,
		cfg.DBTimezone,
	)

	// 2. Buka koneksi menggunakan driver "pgx"
	// Kita menggunakan sqlx.Connect alih-alih sql.Open + sqlx.NewDb
	// karena lebih ringkas dan langsung melakukan Ping.
	db, err := sqlx.Connect("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("gagal terhubung ke database: %w", err)
	}

	// 3. (Opsional) Set koneksi pool
	// db.SetMaxOpenConns(10)
	// db.SetMaxIdleConns(5)
	// db.SetConnMaxLifetime(time.Minute * 10)

	log.Println("Koneksi database berhasil dibuat.")
	return db, nil
}
