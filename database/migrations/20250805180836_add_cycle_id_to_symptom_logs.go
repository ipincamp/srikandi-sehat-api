package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func AddCycleIdToSymptomLogs() *gormigrate.Migration {
	type SymptomLog struct {
		MenstrualCycleID uint `gorm:"nullable;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	}

	return &gormigrate.Migration{
		ID: "20250805180836",

		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&SymptomLog{})
		},

		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropColumn(&SymptomLog{}, "MenstrualCycleID")
		},
	}
}
