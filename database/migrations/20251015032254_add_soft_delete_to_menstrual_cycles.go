package migrations

import (
	"database/sql"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func AddSoftDeleteToMenstrualCycles() *gormigrate.Migration {
	type MenstrualCycle struct {
		DeletionReason sql.NullString `gorm:"type:text"`
		gorm.DeletedAt
	}

	return &gormigrate.Migration{
		ID: "20251015032254",

		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&MenstrualCycle{})
		},

		Rollback: func(tx *gorm.DB) error {
			if err := tx.Migrator().DropColumn(&MenstrualCycle{}, "deleted_at"); err != nil {
				return err
			}
			return tx.Migrator().DropColumn(&MenstrualCycle{}, "deletion_reason")
		},
	}
}
