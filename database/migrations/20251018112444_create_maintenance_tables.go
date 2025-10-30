package migrations

import (
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func CreateMaintenanceTables() *gormigrate.Migration {
	type Setting struct {
		Key       string    `gorm:"primaryKey;type:varchar(100)"`
		Value     string    `gorm:"type:text"`
		CreatedAt time.Time `gorm:"autoCreateTime"`
		UpdatedAt time.Time `gorm:"autoUpdateTime"`
	}

	type MaintenanceWhitelist struct {
		UserID    uint      `gorm:"primaryKey;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
		CreatedAt time.Time `gorm:"autoCreateTime"`
	}

	return &gormigrate.Migration{
		ID: "20251018112444",

		Migrate: func(tx *gorm.DB) error {
			if err := tx.AutoMigrate(&Setting{}); err != nil {
				return err
			}
			if err := tx.AutoMigrate(&MaintenanceWhitelist{}); err != nil {
				return err
			}
			// Initialize default setting
			setting := Setting{Key: "maintenance_mode_active", Value: "false"}
			return tx.FirstOrCreate(&setting).Error
		},

		Rollback: func(tx *gorm.DB) error {
			if err := tx.Migrator().DropTable(&MaintenanceWhitelist{}); err != nil {
				return err
			}
			if err := tx.Migrator().DropTable(&Setting{}); err != nil {
				return err
			}
			return nil
		},
	}
}
