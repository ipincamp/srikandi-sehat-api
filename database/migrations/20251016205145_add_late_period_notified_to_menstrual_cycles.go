package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func AddLatePeriodNotifiedToMenstrualCycles() *gormigrate.Migration {
	type MenstrualCycle struct {
		LatePeriodNotified bool `gorm:"default:false"`
	}

	return &gormigrate.Migration{
		ID: "20251016205145",

		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&MenstrualCycle{})
		},

		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropColumn(&MenstrualCycle{}, "late_period_notified")
		},
	}
}
