package database

import (
	"fmt"
	"ipincamp/srikandi-sehat/config"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDB() {
	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True",
		config.Get("DB_USER"),
		config.Get("DB_PASS"),
		config.Get("DB_HOST"),
		config.Get("DB_PORT"),
		config.Get("DB_NAME"),
	)

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Silent,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
		NowFunc: func() time.Time {
			return time.Now().Local()
		},
	})
	if err != nil {
		log.Fatalf("[DB] Connection failed: %v", err)
	}

	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("[DB] Failed to get basic DB connection: %v", err)
	}

	_, err = sqlDB.Exec("SET time_zone = 'Asia/Jakarta'")
	if err != nil {
		log.Fatalf("[DB] Failed to set database timezone: %v", err)
	}

	log.Println("[DB] Connected to MySQL database")
}
