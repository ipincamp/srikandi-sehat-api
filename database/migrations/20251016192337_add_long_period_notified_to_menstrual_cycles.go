package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func AddLongPeriodNotifiedToMenstrualCycles() *gormigrate.Migration {
	type MenstrualCycle struct {
		LongPeriodNotified bool `gorm:"default:false"`
	}

	return &gormigrate.Migration{
		ID: "20251016192337",

		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&MenstrualCycle{})
		},

		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropColumn(&MenstrualCycle{}, "long_period_notified")
		},
	}
}
