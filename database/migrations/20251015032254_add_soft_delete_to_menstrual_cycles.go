package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func AddSoftDeleteToMenstrualCycles() *gormigrate.Migration {
	type MenstrualCycle struct {
		gorm.DeletedAt
	}

	return &gormigrate.Migration{
		ID: "20251015032254",

		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&MenstrualCycle{})
		},

		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropColumn(&MenstrualCycle{}, "deleted_at")
		},
	}
}
