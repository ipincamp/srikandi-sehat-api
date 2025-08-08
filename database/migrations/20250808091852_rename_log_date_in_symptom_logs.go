package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func RenameLogDateInSymptomLogs() *gormigrate.Migration {
	type SymptomLog struct{}

	return &gormigrate.Migration{
		ID: "20250808091852",

		Migrate: func(tx *gorm.DB) error {
			return tx.Exec("ALTER TABLE symptom_logs CHANGE COLUMN log_date logged_at DATETIME(3)").Error
		},

		Rollback: func(tx *gorm.DB) error {
			return tx.Exec("ALTER TABLE symptom_logs CHANGE COLUMN logged_at log_date DATETIME(3)").Error
		},
	}
}
