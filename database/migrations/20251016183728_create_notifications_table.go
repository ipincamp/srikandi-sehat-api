package migrations

import (
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func CreateNotificationsTable() *gormigrate.Migration {
	type Notification struct {
		ID        uint      `gorm:"primarykey"`
		Title     string    `gorm:"type:varchar(255);not null"`
		Body      string    `gorm:"type:text"`
		IsRead    bool      `gorm:"default:false"`
		UserID    uint      `gorm:"not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
		CreatedAt time.Time `gorm:"autoCreateTime"`
		UpdatedAt time.Time `gorm:"autoUpdateTime"`
	}

	return &gormigrate.Migration{
		ID: "20251016183728",

		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&Notification{})
		},

		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable(&Notification{})
		},
	}
}
