package migrations

import (
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func AddIndexesToSymptomLogs() *gormigrate.Migration {
	type SymptomLog struct {
		UserID  uint      `gorm:"index:idx_user_date"`
		LogDate time.Time `gorm:"index:idx_user_date"`
	}

	return &gormigrate.Migration{
		ID: "20250805124318",

		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&SymptomLog{})
		},

		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropIndex(&SymptomLog{}, "idx_user_date")
		},
	}
}
