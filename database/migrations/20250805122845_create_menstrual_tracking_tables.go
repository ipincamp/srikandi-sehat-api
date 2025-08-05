package migrations

import (
	"database/sql"
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func CreateMenstrualTrackingTables() *gormigrate.Migration {
	type MenstrualCycle struct {
		ID             uint `gorm:"primarykey"`
		StartDate      time.Time
		EndDate        sql.NullTime
		PeriodLength   sql.NullInt16
		CycleLength    sql.NullInt16
		IsPeriodNormal sql.NullBool
		IsCycleNormal  sql.NullBool
		UserID         uint      `gorm:"not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
		CreatedAt      time.Time `gorm:"default:CURRENT_TIMESTAMP"`
		UpdatedAt      time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	}

	type Symptom struct {
		ID        uint      `gorm:"primarykey"`
		Name      string    `gorm:"type:varchar(100);uniqueIndex"`
		Category  string    `gorm:"type:varchar(100)"`
		Type      string    `gorm:"type:enum('BASIC','OPTIONS');default:'BASIC'"`
		CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
		UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	}

	type SymptomOption struct {
		ID        uint      `gorm:"primarykey"`
		Name      string    `gorm:"type:varchar(100)"`
		Value     string    `gorm:"type:varchar(255)"`
		SymptomID uint      `gorm:"not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
		CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
		UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	}

	type SymptomLog struct {
		ID        uint `gorm:"primarykey"`
		LogDate   time.Time
		Note      string    `gorm:"type:text"`
		UserID    uint      `gorm:"not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
		CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
		UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	}

	type SymptomLogDetail struct {
		ID              uint          `gorm:"primarykey"`
		SymptomLogID    uint          `gorm:"not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
		SymptomID       uint          `gorm:"not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
		SymptomOptionID sql.NullInt64 `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	}

	type Recommendation struct {
		ID          uint      `gorm:"primarykey"`
		Title       string    `gorm:"type:varchar(255)"`
		Description string    `gorm:"type:text"`
		Source      string    `gorm:"type:varchar(255)"`
		SymptomID   uint      `gorm:"not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
		CreatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP"`
		UpdatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	}

	return &gormigrate.Migration{
		ID: "20250805122845",

		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(
				&MenstrualCycle{},
				&Symptom{},
				&SymptomOption{},
				&SymptomLog{},
				&SymptomLogDetail{},
				&Recommendation{},
			)
		},

		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable(
				&Recommendation{},
				&SymptomLogDetail{},
				&SymptomLog{},
				&SymptomOption{},
				&Symptom{},
				&MenstrualCycle{},
			)
		},
	}
}
