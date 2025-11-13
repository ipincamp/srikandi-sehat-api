package migrations

import (
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func CreateActivityLogsTable() *gormigrate.Migration {
	type ActivityLog struct {
		ID          int64     `gorm:"primarykey;autoIncrement"`
		ActorUserID *string   `gorm:"type:uuid;index"`
		Action      string    `gorm:"type:varchar(100);not null"`
		TargetTable *string   `gorm:"type:varchar(100);index,priority:1"`
		TargetID    *string   `gorm:"type:uuid;index,priority:2"`
		Changes     *string   `gorm:"type:jsonb"`
		IPAddress   *string   `gorm:"type:varchar(50)"`
		UserAgent   *string   `gorm:"type:text"`
		RequestID   *string   `gorm:"type:uuid"`
		Timestamp   time.Time `gorm:"type:timestamptz;not null;default:NOW()"`
	}

	return &gormigrate.Migration{
		ID: "20251109151006",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&ActivityLog{})
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable(&ActivityLog{})
		},
	}
}
